package handler

import (
	"fmt"
	"insist-backend-golang/internal/dto"
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

type AuthHandler struct {
	authService          *service.AuthService
	passwordResetService *service.PasswordResetService
}

func NewAuthHandler(authService *service.AuthService, passwordResetService *service.PasswordResetService) *AuthHandler {
	return &AuthHandler{
		authService:          authService,
		passwordResetService: passwordResetService,
	}
}

// Login godoc
// @Summary Login a user
// @Description Login with username and password, returns access token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body dto.UserLogin true "Login details"
// @Success 200 {object} map[string]interface{} "Login Successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 401 {object} map[string]interface{} "Unauthorized: Incorrect password, inactive user"
// @Failure 403 {object} map[string]interface{} "Forbidden: Two-factor authentication required"
// @Failure 404 {object} map[string]interface{} "Not Found: User not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var input dto.UserLogin
	if err := c.BodyParser(&input); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	user, err := h.authService.GetByUsername(input.Username)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "User not found"))
	}

	if !user.IsActive {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusUnauthorized, "User is not active"))
	}

	if !pkg.CheckPassword(input.Password, user.Password) {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusUnauthorized, "Invalid password"))
	}

	if user.IsTwoFa {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusForbidden, "Two-factor authentication is required"))
	}

	accessToken, err := pkg.GenerateAccessToken(user.ID)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	refreshToken, err := pkg.GenerateRefreshToken(user.ID)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	if err := h.authService.UpdateRefreshToken(user.ID, refreshToken); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   false,
		SameSite: fiber.CookieSameSiteNoneMode,
		Path:     "/",
		Expires:  time.Now().Add(1 * 24 * time.Hour),
	})

	return pkg.Response(c, fiber.StatusOK, "You have successfully logged in", fiber.Map{
		"access_token": accessToken,
	})
}

// TwoFactorAuth godoc
// @Summary Verify Two-Factor Authentication
// @Description Validate the user's two-factor authentication using an OTP key
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body dto.TwoFactorAuth true "Two-factor authentication details"
// @Success 200 {object} map[string]interface{} "Two-factor authentication successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 401 {object} map[string]interface{} "Unauthorized: Invalid OTP"
// @Failure 404 {object} map[string]interface{} "Not Found: User not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /auth/two-fa [post]
func (h *AuthHandler) TwoFactorAuth(c *fiber.Ctx) error {
	var input dto.TwoFactorAuth
	if err := c.BodyParser(&input); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	user, err := h.authService.GetByUsername(input.Username)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "User not found"))
	}

	valid, err := totp.ValidateCustom(input.OtpKey, user.OtpKey, time.Now(), totp.ValidateOpts{
		Period:    30,
		Skew:      0,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})

	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	if !valid {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusUnauthorized, "Invalid OTP"))
	}

	accessToken, err := pkg.GenerateAccessToken(user.ID)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	refreshToken, err := pkg.GenerateRefreshToken(user.ID)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	if err := h.authService.UpdateRefreshToken(user.ID, refreshToken); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
		Path:     "/",
		Expires:  time.Now().Add(1 * 24 * time.Hour),
	})

	return pkg.Response(c, fiber.StatusOK, "You have successfully logged in", fiber.Map{
		"access_token": accessToken,
	})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Renew the access token using a valid refresh token from cookies
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Access token renewed successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized: Invalid or missing refresh token"
// @Failure 404 {object} map[string]interface{} "Not Found: User not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /auth/token [get]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")

	if refreshToken == "" {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusUnauthorized, "Missing refresh token"))
	}

	userID, err := pkg.VerifyRefreshToken(refreshToken)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token"))
	}

	user, err := h.authService.GetByID(userID)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "User not found"))
	}

	if !user.IsActive {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusUnauthorized, "User is not active"))
	}

	accessToken, err := pkg.GenerateAccessToken(user.ID)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Access token renewed successfully", fiber.Map{
		"access_token": accessToken,
	})
}

// Logout godoc
// @Summary Logout a user
// @Description Log out the user by invalidating the refresh token
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Logout successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized: Invalid or missing refresh token"
// @Failure 404 {object} map[string]interface{} "Not Found: User not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /auth/logout [delete]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")

	if refreshToken == "" {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusUnauthorized, "Missing refresh token"))
	}

	userID, err := pkg.VerifyRefreshToken(refreshToken)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusUnauthorized, "Invalid refresh token"))
	}

	_, err = h.authService.GetByID(userID)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "User not found"))
	}

	if err := h.authService.UpdateRefreshToken(userID, ""); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour),
	})

	return pkg.Response(c, fiber.StatusOK, "Logout successfully", nil)
}

// GetUserInfo godoc
// @Summary Get user information
// @Description Retrieve the information of the currently logged-in user
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "User information retrieved successfully"
// @Failure 401 {object} map[string]interface{} "Unauthorized: Invalid or missing access token"
// @Failure 404 {object} map[string]interface{} "Not Found: User not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /auth/user-info [get]
func (h *AuthHandler) GetUserInfo(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	user, err := h.authService.GetByID(userID)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "User not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "User information", user)
}

// ChangePasswordAuth godoc
// @Summary Change user password auth
// @Description Change the password for the authenticated user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body dto.ChangePasswordAuth true "Current and new password"
// @Success 200 {object} map[string]interface{} "Password changed successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 401 {object} map[string]interface{} "Unauthorized: Invalid current password"
// @Failure 404 {object} map[string]interface{} "Not Found: User not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /auth/change-password [put]
func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var input dto.ChangePasswordAuth
	if err := c.BodyParser(&input); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	user, err := h.authService.GetByID(userID)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "User not found"))
	}

	if !pkg.CheckPassword(input.CurrentPassword, user.Password) {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusUnauthorized, "Invalid password"))
	}

	hashedPassword, err := pkg.HashPassword(input.NewPassword)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	if err := h.authService.UpdatePassword(userID, hashedPassword); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Password changed successfully", nil)
}

// SetTwoFactorAuth godoc
// @Summary Set up two-factor authentication for a user
// @Description Configures two-factor authentication and generates an OTP key for the specified user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{} "Two-factor authentication update successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 404 {object} map[string]interface{} "Not Found: User not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /auth/{id}/two-fa [put]
func (h *AuthHandler) SetTwoFactorAuth(c *fiber.Ctx) error {
	userID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	user, err := h.authService.GetByID(uint(userID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "User not found"))
	}

	if user.Email == "" {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, "Email not found"))
	}

	err = h.authService.UpdateTwoFactorAuth(uint(userID), true)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	if user.OtpKey != "" {
		result := map[string]interface{}{
			"otp_key": user.OtpKey,
			"otp_url": user.OtpUrl,
		}

		return pkg.Response(c, fiber.StatusOK, "Two-factor authentication already set", result)
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "INSIST",
		AccountName: strings.ToUpper(user.Name),
	})
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	user.OtpUrl = key.URL()
	qrCode, err := qrcode.New(key.URL(), qrcode.Medium)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	fmt.Println(user.CreatedAt)
	folderPath := "./public/images/qrcodeotp"
	fileName := fmt.Sprintf("%s_%s.png", user.CreatedAt.Format("20060102"), key.Secret())

	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
		}
	}

	fullFilePath := filepath.Join(folderPath, fileName)
	err = qrCode.WriteFile(256, fullFilePath)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	err = h.authService.UpdateTwoFactorAuthKey(uint(userID), key.Secret(), key.URL())
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"otp_key": key.Secret(),
		"otp_url": key.URL(),
	}

	return pkg.Response(c, fiber.StatusOK, "Two-factor authentication update successfully", result)
}

// SendPasswordReset godoc
// @Summary Send a password reset link to the user's email
// @Description Sends a password reset link to the specified user after generating a reset token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body model.PasswordReset true "Password reset details"
// @Success 200 {object} map[string]interface{} "Email sent successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 404 {object} map[string]interface{} "Not Found: User not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error: Failed to send email"
// @Router /auth/send-password-reset [post]
func (h *AuthHandler) SendPasswordReset(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	user, err := h.authService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "User not found"))
	}

	token, err := pkg.GenerateToken()
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	var passwordReset model.PasswordReset
	expirationTime := time.Now().Add(time.Hour * 24)
	passwordReset.IDUser = ID
	passwordReset.Token = token
	passwordReset.ExpiredAt = expirationTime
	passwordReset.IDCreatedby = userID
	passwordReset.IDUpdatedby = userID

	err = h.passwordResetService.CreatePasswordReset(&passwordReset)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	template := pkg.TemplateSendOTP(user.Name, user.Username, user.OtpKey, token)

	emailSender := pkg.NewEmailSender(os.Getenv("MAIL_SMTP"), 587, os.Getenv("MAIL_EMAIL"), os.Getenv("MAIL_PASSWORD"))
	err = emailSender.SendEmail(user.Email, "Reset Password & Aktivasi 2FA", "text/html", template)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, "Email sending failed"))
	}

	return pkg.Response(c, fiber.StatusOK, "Email sent successfully", nil)
}

// PasswordReset godoc
// @Summary Reset user password
// @Description Resets the user's password using a valid, unused, and non-expired token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param token query string true "Password reset token"
// @Param input body dto.ChangePassword true "New password details"
// @Success 200 {object} map[string]interface{} "Password reset successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Passwords do not match or input is invalid"
// @Failure 401 {object} map[string]interface{} "Unauthorized: Expired or used token"
// @Failure 404 {object} map[string]interface{} "Not Found: Invalid token"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /auth/password-reset [post]
func (h *AuthHandler) PasswordReset(c *fiber.Ctx) error {
	token := c.Query("token", "")

	passwordReset, err := h.passwordResetService.GetPasswordResetByToken(token)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Invalid token"))
	}

	if passwordReset.ExpiredAt.Before(time.Now()) {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusUnauthorized, "Expired token"))
	}

	if passwordReset.IsUsed {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusUnauthorized, "Token has already been used"))
	}

	var input dto.ChangePassword
	if err := c.BodyParser(&input); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	if input.Password != input.ConfirmPassword {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, "Password do not match"))
	}

	hashedPassword, err := pkg.HashPassword(input.Password)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	err = h.authService.UpdatePassword(uint(passwordReset.IDUser), hashedPassword)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	passwordReset.IsUsed = true
	err = h.passwordResetService.UpdateUsed(passwordReset.ID, true)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Password Reset Successfully", nil)
}
