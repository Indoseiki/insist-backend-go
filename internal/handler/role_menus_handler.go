package handler

import (
	"insist-backend-golang/internal/dto"
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type RoleMenuHandler struct {
	roleMenuService *service.RoleMenuService
}

func NewRoleMenuHandler(roleMenuService *service.RoleMenuService) *RoleMenuHandler {
	return &RoleMenuHandler{roleMenuService: roleMenuService}
}

// GetRoleMenus godoc
// @Summary Get a list of Role Menus
// @Description Retrieves Role Menus with pagination and optional search
// @Tags Role Menu
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering Role Menu"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/role-menu [get]
func (h *RoleMenuHandler) GetRoleMenus(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search")
	offset := (page - 1) * rows

	total, err := h.roleMenuService.GetTotal(search)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	roleMenus, err := h.roleMenuService.GetAll(offset, rows, search)
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
		"items": roleMenus,
		"pagination": map[string]interface{}{
			"current_page":  page,
			"next_page":     nextPage,
			"total_pages":   totalPages,
			"rows_per_page": rows,
			"total_rows":    total,
		},
	}

	if len(roleMenus) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetRoleMenu godoc
// @Summary Get Role Menu by ID
// @Description Retrieve a specific Role Menu by its ID
// @Tags Role Menu
// @Accept json
// @Produce json
// @Param id path int true "Role Menu ID"
// @Success 200 {object} map[string]interface{} "Role Menu found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Role Menu not found"
// @Router /admin/master/role-menu/{id} [get]
func (h *RoleMenuHandler) GetRoleMenu(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	roleMenu, err := h.roleMenuService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Role Menu not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Role Menu found successfully", roleMenu)
}

// UpdateRoleMenu godoc
// @Summary Update a new Role Menu
// @Description Update a new Role Menu with the provided details
// @Tags Role Menu
// @Accept json
// @Produce json
// @Param RoleMenu body dto.RoleMenus true "Role Menu details"
// @Success 201 {object} map[string]interface{} "Role Menu created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/role-menu [post]
func (h *RoleMenuHandler) UpdateRoleMenu(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var input dto.RoleMenus
	if err := c.BodyParser(&input); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var roleMenus []model.MstRoleMenu
	for _, menuID := range input.IDMenu {
		roleMenus = append(roleMenus, model.MstRoleMenu{
			IDMenu: menuID,
			IDRole: uint(ID),
		})
	}

	err = h.roleMenuService.Delete(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	err = h.roleMenuService.Create(&roleMenus)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusCreated, "Role Menu created successfully", roleMenus)
}
