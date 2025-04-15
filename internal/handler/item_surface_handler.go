package handler

import (
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type ItemSurfaceHandler struct {
	itemSurfaceService *service.ItemSurfaceService
}

func NewItemSurfaceHandler(itemSurfaceService *service.ItemSurfaceService) *ItemSurfaceHandler {
	return &ItemSurfaceHandler{itemSurfaceService: itemSurfaceService}
}

// GetItemSurfaces godoc
// @Summary Get a list of Item Surfaces
// @Description Retrieves Item Surfaces with pagination, search, filtering, and sorting options
// @Tags Item Surface
// @Accept json
// @Produce json
// @Param categoryCode query string false "Filter by Item Category Code"
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering Item Surface"
// @Param sortBy query string false "Field to sort by"
// @Param sortDirection query boolean false "Sort direction: true for ascending, false for descending"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /general/master/item/surface [get]
func (h *ItemSurfaceHandler) GetItemSurfaces(c *fiber.Ctx) error {
	categoryCode := c.Query("categoryCode", "")
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	total, err := h.itemSurfaceService.GetTotal(search, categoryCode)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	itemSurfaces, err := h.itemSurfaceService.GetAll(offset, rows, search, sortBy, sortDirection, categoryCode)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	totalPages := int(math.Ceil(float64(total) / float64(rows)))

	var start, end *int
	if total > 0 {
		startVal := offset + 1
		endVal := int(math.Min(float64(offset+rows), float64(total)))
		start = &startVal
		end = &endVal
	}

	var nextPage *int
	if page < totalPages {
		next := page + 1
		nextPage = &next
	}

	result := map[string]interface{}{
		"items": itemSurfaces,
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

	if len(itemSurfaces) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetItemSurface godoc
// @Summary Get Item Surface by ID
// @Description Retrieve a specific Item Surface by its ID
// @Tags Item Surface
// @Accept json
// @Produce json
// @Param id path int true "Item Surface ID"
// @Success 200 {object} map[string]interface{} "Item Surface found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Item Surface not found"
// @Router /general/master/item/surface/{id} [get]
func (h *ItemSurfaceHandler) GetItemSurface(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	itemSurface, err := h.itemSurfaceService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Item Surface not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Item Surface found successfully", itemSurface)
}

// CreateItemSurface godoc
// @Summary Create a new Item Surface
// @Description Create a new Item Surface with the provided details
// @Tags Item Surface
// @Accept json
// @Produce json
// @Param ItemSurface body model.MstItemSurface true "Item Surface details"
// @Success 201 {object} map[string]interface{} "Item Surface created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /general/master/item/surface [post]
func (h *ItemSurfaceHandler) CreateItemSurface(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var itemSurface model.MstItemSurface
	if err := c.BodyParser(&itemSurface); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	itemSurface.IDCreatedby = userID
	itemSurface.IDUpdatedby = userID

	if err := h.itemSurfaceService.Create(&itemSurface); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusCreated, "Item Surface created successfully", map[string]interface{}{
		"id": itemSurface.ID,
	})
}

// UpdateItemSurface godoc
// @Summary Update an existing Item Surface
// @Description Update the details of an existing Item Surface by its ID
// @Tags Item Surface
// @Accept json
// @Produce json
// @Param id path int true "Item Surface ID"
// @Param ItemSurface body model.MstItemSurface true "Updated Item Surface details"
// @Success 200 {object} map[string]interface{} "Item Surface updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 404 {object} map[string]interface{} "Not Found: Item Surface not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /general/master/item/surface/{id} [put]
func (h *ItemSurfaceHandler) UpdateItemSurface(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	itemSurface, err := h.itemSurfaceService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Item Surface not found"))
	}

	if err := c.BodyParser(itemSurface); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	itemSurface.IDUpdatedby = userID

	if err := h.itemSurfaceService.Update(itemSurface); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Item Surface updated successfully", map[string]interface{}{
		"id": itemSurface.ID,
	})
}

// DeleteItemSurface godoc
// @Summary Delete an Item Surface
// @Description Delete an Item Surface by its ID
// @Tags Item Surface
// @Param id path int true "Item Surface ID"
// @Success 200 {object} map[string]interface{} "Item Surface deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Item Surface not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /general/master/item/surface/{id} [delete]
func (h *ItemSurfaceHandler) DeleteItemSurface(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	itemSurface, err := h.itemSurfaceService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Item Surface not found"))
	}

	if err := h.itemSurfaceService.Delete(itemSurface); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Item Surface deleted successfully", nil)
}
