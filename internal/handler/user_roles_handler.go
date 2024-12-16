package handler

import (
	"insist-backend-golang/internal/dto"
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type UserRoleHandler struct {
	userRoleService *service.UserRoleService
}

func NewUserRoleHandler(userRoleService *service.UserRoleService) *UserRoleHandler {
	return &UserRoleHandler{userRoleService: userRoleService}
}

// GetUserRoles godoc
// @Summary Get a list of User Roles
// @Description Retrieves User Roles with pagination and optional search
// @Tags User Role
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering User Role"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/user-role [get]
func (h *UserRoleHandler) GetUserRoles(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search")
	offset := (page - 1) * rows

	total, err := h.userRoleService.GetTotal(search)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	userRoles, err := h.userRoleService.GetAll(offset, rows, search)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	totalPages := int(math.Ceil(float64(total) / float64(rows)))
	var nextPage *int
	if page < totalPages {
		nextPageVal := page + 1
		nextPage = &nextPageVal
	}

	result := map[string]interface{}{
		"items": userRoles,
		"pagination": map[string]interface{}{
			"current_page":  page,
			"next_page":     nextPage,
			"total_pages":   totalPages,
			"rows_per_page": rows,
			"total_rows":    total,
		},
	}

	if len(userRoles) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetUserRole godoc
// @Summary Get User Role by ID
// @Description Retrieve a specific User Role by its ID
// @Tags User Role
// @Accept json
// @Produce json
// @Param id path int true "User Role ID"
// @Success 200 {object} map[string]interface{} "User Role found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: User Role not found"
// @Router /admin/master/user-role/{id} [get]
func (h *UserRoleHandler) GetUserRole(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	userRole, err := h.userRoleService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "User Role not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "User Role found successfully", userRole)
}

// CreateUserRole godoc
// @Summary Create a new User Role
// @Description Create a new User Role with the provided details
// @Tags User Role
// @Accept json
// @Produce json
// @Param UserRole body dto.UserRoles true "User Role details"
// @Success 201 {object} map[string]interface{} "User Role created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/user-role [post]
func (h *UserRoleHandler) UpdateUserRole(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var input dto.UserRoles
	if err := c.BodyParser(&input); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var userRoles []model.MstUserRole
	for _, roleID := range input.IDRole {
		userRoles = append(userRoles, model.MstUserRole{
			IDUser: uint(ID),
			IDRole: roleID,
		})
	}

	err = h.userRoleService.Delete(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	err = h.userRoleService.Create(&userRoles)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusCreated, "User Role created successfully", userRoles)
}
