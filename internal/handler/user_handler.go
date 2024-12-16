package handler

import (
	"insist-backend-golang/internal/dto"
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetUsers godoc
// @Summary Get a list of users
// @Description Retrieves a paginated list of users based on search criteria
// @Tags Users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of users per page" default(20)
// @Param search query string false "Search keyword for filtering users"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/users [get]
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search")
	idDept := c.QueryInt("idDept", 0)
	isActive := c.Query("isActive", "")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	total, err := h.userService.GetTotal(search, uint(idDept), isActive)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	users, err := h.userService.GetAll(offset, rows, search, uint(idDept), isActive, sortBy, sortDirection)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	totalPages := int(math.Ceil(float64(total) / float64(rows)))

	var start *int
	if int(total) == 0 {
		start = nil
	} else {
		value := offset + 1
		start = &value
	}

	var end *int
	if int(total) == 0 {
		end = nil
	} else {
		value := int(math.Min(float64(offset+rows), float64(total)))
		end = &value
	}

	var nextPage *int
	if page < totalPages {
		nextPageVal := page + 1
		nextPage = &nextPageVal
	}

	result := map[string]interface{}{
		"items": users,
		"pagination": map[string]interface{}{
			"current_page":  page,
			"next_page":     nextPage,
			"total_pages":   totalPages,
			"rows_per_page": rows,
			"total_rows":    total,
			"from":          start,
			"to":            end,
		},
	}

	if len(users) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Retrieves user details based on the provided user ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{} "User found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: User not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/users/{id} [get]
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	users, err := h.userService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "User not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "User found successfully", users)
}

// CreateUser godoc
// @Summary Create a new user
// @Description Creates a new user with the provided details
// @Tags Users
// @Accept json
// @Produce json
// @Param user body model.MstUser true "User data"
// @Success 201 {object} map[string]interface{} "User created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/users [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var user model.MstUser
	if err := c.BodyParser(&user); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	hashedPassword, err := pkg.HashPassword(user.Password)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	user.Password = string(hashedPassword)
	user.IDCreatedby = userID
	user.IDUpdatedby = userID

	err = h.userService.Create(&user)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": user.ID,
	}

	return pkg.Response(c, fiber.StatusCreated, "User created successfully", result)
}

// UpdateUser godoc
// @Summary Update an existing user
// @Description Updates user details based on the provided user ID and new data
// @Tags Users
// @Accept json
// @Produce json
// @Param id query int true "User ID"
// @Param user body model.MstUser true "Updated user data"
// @Success 200 {object} map[string]interface{} "User updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 404 {object} map[string]interface{} "Not Found: User not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/users/{id} [put]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, "Invalid user ID"))
	}

	_, err = h.userService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "User not found"))
	}

	var updatedUser model.MstUser
	if err := c.BodyParser(&updatedUser); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, "Invalid request body"))
	}

	updatedUser.ID = uint(ID)
	updatedUser.IDUpdatedby = userID

	err = h.userService.Update(&updatedUser)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, "Failed to update user"))
	}

	result := map[string]interface{}{
		"id": updatedUser.ID,
	}

	return pkg.Response(c, fiber.StatusOK, "User updated successfully", result)
}

// DeleteUser godoc
// @Summary Delete a user by ID
// @Description Deletes an existing user based on the provided user ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{} "User deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid user ID"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	user, err := h.userService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "User not found"))
	}

	err = h.userService.Delete(user)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "User deleted successfully", nil)
}

// ChangePassword godoc
// @Summary Change user password user
// @Description Changes the password for user
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param input body dto.ChangePassword true "New password and confirmation password"
// @Success 200 {object} map[string]interface{} "Password changed successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input or Confirm password does not match"
// @Failure 404 {object} map[string]interface{} "Not Found: User not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/users/{id}/change-password [put]
func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var input dto.ChangePassword
	if err := c.BodyParser(&input); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	_, err = h.userService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "User not found"))
	}

	if input.Password != input.ConfirmPassword {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, "Confirm password does not match"))
	}

	hashedPassword, err := pkg.HashPassword(input.Password)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	err = h.userService.UpdatePassword(uint(ID), hashedPassword)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Password changed successfully", nil)
}
