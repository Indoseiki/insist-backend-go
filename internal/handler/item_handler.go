package handler

import (
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type ItemHandler struct {
	itemService *service.ItemService
}

func NewItemHandler(itemService *service.ItemService) *ItemHandler {
	return &ItemHandler{itemService: itemService}
}

// GetItems godoc
// @Summary Get list of Items
// @Tags Item
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param rows query int false "Rows per page"
// @Param search query string false "Search"
// @Param idItemCategory query string false "Item Category ID"
// @Param sortBy query string false "Sort by field"
// @Param sortDirection query boolean false "true = ASC, false = DESC"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /general/master/items [get]
func (h *ItemHandler) GetItems(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search", "")
	idItemCategory := c.Query("idItemCategory", "")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection", true)
	offset := (page - 1) * rows

	total, err := h.itemService.GetTotal(search, idItemCategory)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	items, err := h.itemService.GetAll(offset, rows, search, sortBy, sortDirection, idItemCategory)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	totalPages := int(math.Ceil(float64(total) / float64(rows)))
	start := offset + 1
	end := int(math.Min(float64(offset+rows), float64(total)))

	result := map[string]interface{}{
		"items": items,
		"pagination": map[string]interface{}{
			"current_page":  page,
			"total_pages":   totalPages,
			"rows_per_page": rows,
			"total_rows":    total,
			"from":          start,
			"to":            end,
		},
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetItem godoc
// @Summary Get item by ID
// @Tags Item
// @Param id path int true "Item ID"
// @Success 200 {object} model.MstItem
// @Failure 404 {object} map[string]interface{}
// @Router /general/master/items/{id} [get]
func (h *ItemHandler) GetItem(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, "Invalid ID"))
	}

	item, err := h.itemService.GetByID(uint(id))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Item not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Item found successfully", item)
}

// CreateItem godoc
// @Summary Create new item
// @Tags Item
// @Accept json
// @Produce json
// @Param item body model.MstItem true "Item Body"
// @Success 201 {object} map[string]interface{}
// @Router /general/master/items [post]
func (h *ItemHandler) CreateItem(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var item model.MstItem
	if err := c.BodyParser(&item); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	item.IDCreatedby = userID
	item.IDUpdatedby = userID

	if err := h.itemService.Create(&item); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusCreated, "Item created successfully", map[string]interface{}{"id": item.ID})
}

// UpdateItem godoc
// @Summary Update item
// @Tags Item
// @Param id path int true "Item ID"
// @Param item body model.MstItem true "Updated Item"
// @Success 200 {object} map[string]interface{}
// @Router /general/master/items/{id} [put]
func (h *ItemHandler) UpdateItem(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, "Invalid ID"))
	}

	item, err := h.itemService.GetByID(uint(id))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Item not found"))
	}

	if err := c.BodyParser(item); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	item.IDUpdatedby = userID

	if err := h.itemService.Update(item); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Item updated successfully", map[string]interface{}{"id": item.ID})
}

// DeleteItem godoc
// @Summary Delete item
// @Tags Item
// @Param id path int true "Item ID"
// @Success 200 {object} map[string]interface{}
// @Router /general/master/items/{id} [delete]
func (h *ItemHandler) DeleteItem(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, "Invalid ID"))
	}

	item, err := h.itemService.GetByID(uint(id))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Item not found"))
	}

	if err := h.itemService.Delete(item); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Item deleted successfully", nil)
}
