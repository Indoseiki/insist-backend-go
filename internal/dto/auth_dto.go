package dto

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TwoFactorAuth struct {
	OtpKey   string `json:"otp_key"`
	Username string `json:"username"`
}

type OTPKey struct {
	OtpKey string `json:"otp_key"`
}

type ChangePasswordAuth struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}
