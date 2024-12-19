package handler

import (
	"insist-backend-golang/internal/dto"
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"

	"github.com/gofiber/fiber/v2"
)

type RolePermissionHandler struct {
	rolePermissionService *service.RolePermissionService
}

func NewRolePermissionHandler(rolePermissionService *service.RolePermissionService) *RolePermissionHandler {
	return &RolePermissionHandler{rolePermissionService: rolePermissionService}
}

// GetMenuTreeByRole godoc
// @Summary Get menu tree by role
// @Description Retrieve the menu tree associated with a specific role by role ID
// @Tags Role Permission
// @Produce json
// @Param id path int true "Role ID"
// @Success 200 {object} map[string]interface{} "Menu tree role permission retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/role-permission/{id} [get]
func (h *RolePermissionHandler) GetMenuTreeByRole(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	tree, err := h.rolePermissionService.GetMenuTreeByRole(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Menu tree role permission retrieved successfully", tree)
}

func (h *RolePermissionHandler) GetMenuByPath(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	path := c.Query("path")

	menus, err := h.rolePermissionService.GetByPath(userID, path)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Menu not found"))
	}

	if len(menus) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Menu not found"))
	}

	var result dto.MenuWithPermissions

	for i, menu := range menus {
		if i == 0 {
			result = menu
		} else {
			result.IsCreate = result.IsCreate || menu.IsCreate
			result.IsUpdate = result.IsUpdate || menu.IsUpdate
			result.IsDelete = result.IsDelete || menu.IsDelete
		}
	}

	return pkg.Response(c, fiber.StatusOK, "Menu found successfully", result)
}

// CreateRolePermission godoc
// @Summary Create or update a role permission
// @Description Create a new role permission or update an existing one if it already exists
// @Tags Role Permission
// @Accept json
// @Produce json
// @Param rolePermission body model.MstRolePermission true "Role Permission data"
// @Success 201 {object} map[string]interface{} "Role permission created or updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/role-permission [post]
func (h *RolePermissionHandler) UpdateOrCreateRolePermission(c *fiber.Ctx) error {
	var rolePermission model.MstRolePermission
	if err := c.BodyParser(&rolePermission); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	err := h.rolePermissionService.UpdateOrCreateRolePermission(&rolePermission)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusCreated, "Role Permission updated successfully", rolePermission)
}
