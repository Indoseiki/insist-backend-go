package handler

import (
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type ItemSourceHandler struct {
	itemSourceService *service.ItemSourceService
}

func NewItemSourceHandler(svc *service.ItemSourceService) *ItemSourceHandler {
	return &ItemSourceHandler{itemSourceService: svc}
}

// GetItemSources godoc
// @Summary Get a list of Item Sources
// @Description Retrieves Item Sources with pagination, search, filtering, and sorting options
// @Tags Item Source
// @Accept json
// @Produce json
// @Param categoryCode query string false "Filter by Item Category Code"
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering Item Source"
// @Param sortBy query string false "Field to sort by"
// @Param sortDirection query boolean false "Sort direction: true for ascending, false for descending"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /general/master/item/source [get]
func (h *ItemSourceHandler) GetItemSources(c *fiber.Ctx) error {
	categoryCode := c.Query("categoryCode", "")
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	total, err := h.itemSourceService.GetTotal(search, categoryCode)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	sources, err := h.itemSourceService.GetAll(offset, rows, search, sortBy, sortDirection, categoryCode)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	totalPages := int(math.Ceil(float64(total) / float64(rows)))

	var start, end *int
	if int(total) > 0 {
		s := offset + 1
		e := int(math.Min(float64(offset+rows), float64(total)))
		start = &s
		end = &e
	}

	var nextPage *int
	if page < totalPages {
		n := page + 1
		nextPage = &n
	}

	result := map[string]interface{}{
		"items": sources,
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

	if len(sources) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetItemSource godoc
// @Summary Get Item Source by ID
// @Tags Item Source
// @Param id path int true "Item Source ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /general/master/item/source/{id} [get]
func (h *ItemSourceHandler) GetItemSource(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	data, err := h.itemSourceService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Item Source not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Item Source found successfully", data)
}

// CreateItemSource godoc
// @Summary Create a new Item Source
// @Tags Item Source
// @Accept json
// @Produce json
// @Param ItemSource body model.MstItemSource true "Item Source data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /general/master/item/source [post]
func (h *ItemSourceHandler) CreateItemSource(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var data model.MstItemSource
	if err := c.BodyParser(&data); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	data.IDCreatedby = userID
	data.IDUpdatedby = userID

	err := h.itemSourceService.Create(&data)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusCreated, "Item Source created successfully", map[string]interface{}{"id": data.ID})
}

// UpdateItemSource godoc
// @Summary Update Item Source by ID
// @Tags Item Source
// @Accept json
// @Param id path int true "Item Source ID"
// @Param ItemSource body model.MstItemSource true "Updated Item Source data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /general/master/item/source/{id} [put]
func (h *ItemSourceHandler) UpdateItemSource(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	data, err := h.itemSourceService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Item Source not found"))
	}

	if err := c.BodyParser(data); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	data.IDUpdatedby = userID

	err = h.itemSourceService.Update(data)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Item Source updated successfully", map[string]interface{}{"id": data.ID})
}

// DeleteItemSource godoc
// @Summary Delete Item Source by ID
// @Tags Item Source
// @Param id path int true "Item Source ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /general/master/item/source/{id} [delete]
func (h *ItemSourceHandler) DeleteItemSource(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	data, err := h.itemSourceService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Item Source not found"))
	}

	err = h.itemSourceService.Delete(data)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Item Source deleted successfully", nil)
}
