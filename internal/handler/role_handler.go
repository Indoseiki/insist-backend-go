package handler

import (
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type RoleHandler struct {
	roleService *service.RoleService
}

func NewRoleHandler(roleService *service.RoleService) *RoleHandler {
	return &RoleHandler{roleService: roleService}
}

// GetRoles godoc
// @Summary Get a list of roles
// @Description Retrieves roles with pagination and optional search
// @Tags Role
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering role"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/role [get]
func (h *RoleHandler) GetRoles(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	total, err := h.roleService.GetTotal(search)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	roles, err := h.roleService.GetAll(offset, rows, search, sortBy, sortDirection)
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
		"items": roles,
		"pagination": map[string]interface{}{
			"current_page":  page,
			"next_page":     nextPage,
			"total_pages":   totalPages,
			"rows_per_page": rows,
			"total_rows":    total,
		},
	}

	if len(roles) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetRole godoc
// @Summary Get role by ID
// @Description Retrieve a specific role by its ID
// @Tags Role
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} map[string]interface{} "Role found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Role not found"
// @Router /admin/master/role/{id} [get]
func (h *RoleHandler) GetRole(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	role, err := h.roleService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Role not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Role found successfully", role)
}

// CreateRole godoc
// @Summary Create a new role
// @Description Create a new role with the provided details
// @Tags Role
// @Accept json
// @Produce json
// @Param role body model.MstRole true "Role details"
// @Success 201 {object} map[string]interface{} "Role created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/role [post]
func (h *RoleHandler) CreateRole(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var role model.MstRole
	if err := c.BodyParser(&role); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	role.IDCreatedby = userID
	role.IDUpdatedby = userID

	err := h.roleService.Create(&role)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": role.ID,
	}

	return pkg.Response(c, fiber.StatusCreated, "Role created successfully", result)
}

// UpdateRole godoc
// @Summary Update an existing role
// @Description Update the details of an existing role by its ID
// @Tags Role
// @Accept json
// @Produce json
// @Param id path int true "Role ID"
// @Param role body model.MstRole true "Updated role details"
// @Success 200 {object} map[string]interface{} "Role updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 404 {object} map[string]interface{} "Not Found: Role not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/role/{id} [put]
func (h *RoleHandler) UpdateRole(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var role *model.MstRole
	role, err = h.roleService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Role not found"))
	}

	if err := c.BodyParser(role); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	role.ID = uint(ID)
	role.IDUpdatedby = userID

	err = h.roleService.Update(role)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": role.ID,
	}

	return pkg.Response(c, fiber.StatusOK, "Role updated successfully", result)
}

// DeleteRole godoc
// @Summary Delete a role
// @Description Delete a role by its ID
// @Tags Role
// @Param id path int true "Role ID"
// @Success 200 {object} map[string]interface{} "Role deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Role not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/role/{id} [delete]
func (h *RoleHandler) DeleteRole(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	user, err := h.roleService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Role not found"))
	}

	err = h.roleService.Delete(user)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Role deleted successfully", nil)
}
