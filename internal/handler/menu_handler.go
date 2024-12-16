package handler

import (
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type MenuHandler struct {
	menuService *service.MenuService
}

func NewMenuHandler(menuService *service.MenuService) *MenuHandler {
	return &MenuHandler{menuService: menuService}
}

// GetMenus godoc
// @Summary Get a list of menus
// @Description Retrieves menus with pagination and optional search
// @Tags Menu
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering menu"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/menu [get]
func (h *MenuHandler) GetMenus(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	total, err := h.menuService.GetTotal(search)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	menus, err := h.menuService.GetAll(offset, rows, search, sortBy, sortDirection)
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
		"items": menus,
		"pagination": map[string]interface{}{
			"current_page":  page,
			"next_page":     nextPage,
			"total_pages":   totalPages,
			"rows_per_page": rows,
			"total_rows":    total,
		},
	}

	if len(menus) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetMenu godoc
// @Summary Get menu by ID
// @Description Retrieve a specific menu by its ID
// @Tags Menu
// @Accept json
// @Produce json
// @Param id path int true "Menu ID"
// @Success 200 {object} map[string]interface{} "Menu found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Menu not found"
// @Router /admin/master/menu/{id} [get]
func (h *MenuHandler) GetMenu(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	menu, err := h.menuService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Menu not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Menu found successfully", menu)
}

// CreateMenu godoc
// @Summary Create a new menu
// @Description Create a new menu with the provided details
// @Tags Menu
// @Accept json
// @Produce json
// @Param menu body model.MstMenu true "Menu details"
// @Success 201 {object} map[string]interface{} "Menu created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/menu [post]
func (h *MenuHandler) CreateMenu(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var menu model.MstMenu
	if err := c.BodyParser(&menu); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	menu.IDCreatedby = userID
	menu.IDUpdatedby = userID

	err := h.menuService.Create(&menu)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": menu.ID,
	}

	return pkg.Response(c, fiber.StatusCreated, "Menu created successfully", result)
}

// UpdateMenu godoc
// @Summary Update an existing menu
// @Description Update the details of an existing menu by its ID
// @Tags Menu
// @Accept json
// @Produce json
// @Param id path int true "Menu ID"
// @Param menu body model.MstMenu true "Updated menu details"
// @Success 200 {object} map[string]interface{} "Menu updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 404 {object} map[string]interface{} "Not Found: Menu not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/menu/{id} [put]
func (h *MenuHandler) UpdateMenu(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var menu *model.MstMenu
	menu, err = h.menuService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Menu not found"))
	}

	if err := c.BodyParser(menu); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	menu.ID = uint(ID)
	menu.IDUpdatedby = userID

	err = h.menuService.Update(menu)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": menu.ID,
	}

	return pkg.Response(c, fiber.StatusOK, "Menu updated successfully", result)
}

// DeleteMenu godoc
// @Summary Delete a menu
// @Description Delete a menu by its ID
// @Tags Menu
// @Param id path int true "Menu ID"
// @Success 200 {object} map[string]interface{} "Menu deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Menu not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/menu/{id} [delete]
func (h *MenuHandler) DeleteMenu(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	user, err := h.menuService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Menu not found"))
	}

	err = h.menuService.Delete(user)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Menu deleted successfully", nil)
}

// GetMenuTree godoc
// @Summary Get the menu tree
// @Description Retrieve the menu tree structure
// @Tags Menu
// @Success 200 {object} map[string]interface{} "Menu tree retrieved successfully"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/tree-menu [get]
func (h *MenuHandler) GetMenuTree(c *fiber.Ctx) error {
	tree, err := h.menuService.GetMenuTree()
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Menu tree retrieved successfully", tree)
}

// GetMenuTreeByUser godoc
// @Summary Get the menu tree for a user login
// @Description Retrieve the menu tree structure specific to a user login based on their ID
// @Tags Menu
// @Param userID path int true "User ID"
// @Success 200 {object} map[string]interface{} "Menu tree for user retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/tree-menu/user [get]
func (h *MenuHandler) GetMenuTreeByUser(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	tree, err := h.menuService.GetMenuTreeByUser(userID)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Menu tree user retrieved successfully", tree)
}
