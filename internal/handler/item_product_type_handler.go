package handler

import (
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type ItemProductTypeHandler struct {
	itemProductService *service.ItemProductTypeService
}

func NewItemProductTypeHandler(itemProductService *service.ItemProductTypeService) *ItemProductTypeHandler {
	return &ItemProductTypeHandler{itemProductService: itemProductService}
}

// GetItemProductTypes godoc
// @Summary Get a list of Item Product Types
// @Description Retrieves Item Product Types with pagination, search, filtering, and sorting options
// @Tags Item Product Type
// @Accept json
// @Produce json
// @Param idItemProduct query int false "Filter by Item Product ID"
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering Item Product Type"
// @Param sortBy query string false "Field to sort by"
// @Param sortDirection query boolean false "Sort direction: true for ascending, false for descending"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /general/master/item/product-type [get]
func (h *ItemProductTypeHandler) GetItemProductTypes(c *fiber.Ctx) error {
	idItemProduct := c.QueryInt("idItemProduct", 0)
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	total, err := h.itemProductService.GetTotal(search, uint(idItemProduct))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	itemProductTypes, err := h.itemProductService.GetAll(offset, rows, search, sortBy, sortDirection, uint(idItemProduct))
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
		"items": itemProductTypes,
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

	if len(itemProductTypes) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetItemProductType godoc
// @Summary Get Item Product Type by ID
// @Description Retrieve a specific Item Product Type by its ID
// @Tags Item Product Type
// @Accept json
// @Produce json
// @Param id path int true "Item Product Type ID"
// @Success 200 {object} map[string]interface{} "Item Product Type found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Item Product Type not found"
// @Router /general/master/item/product-type/{id} [get]
func (h *ItemProductTypeHandler) GetItemProductType(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	itemProduct, err := h.itemProductService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Item Product Type not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Item Product Type found successfully", itemProduct)
}

// CreateItemProductType godoc
// @Summary Create a new Item Product Type
// @Description Create a new Item Product Type with the provided details
// @Tags Item Product Type
// @Accept json
// @Produce json
// @Param ItemProductType body model.MstItemProductType true "Item Product Type details"
// @Success 201 {object} map[string]interface{} "Item Product Type created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /general/master/item/product-type [post]
func (h *ItemProductTypeHandler) CreateItemProductType(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var itemProduct model.MstItemProductType
	if err := c.BodyParser(&itemProduct); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	itemProduct.IDCreatedby = userID
	itemProduct.IDUpdatedby = userID

	err := h.itemProductService.Create(&itemProduct)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": itemProduct.ID,
	}

	return pkg.Response(c, fiber.StatusCreated, "Item Product Type created successfully", result)
}

// UpdateItemProductType godoc
// @Summary Update an existing Item Product Type
// @Description Update the details of an existing Item Product Type by its ID
// @Tags Item Product Type
// @Accept json
// @Produce json
// @Param id path int true "Item Product Type ID"
// @Param ItemProductType body model.MstItemProductType true "Updated Item Product Type details"
// @Success 200 {object} map[string]interface{} "Item Product Type updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 404 {object} map[string]interface{} "Not Found: Item Product Type not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /general/master/item/product-type/{id} [put]
func (h *ItemProductTypeHandler) UpdateItemProductType(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var itemProduct *model.MstItemProductType
	itemProduct, err = h.itemProductService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Item Product Type not found"))
	}

	if err := c.BodyParser(itemProduct); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	itemProduct.ID = uint(ID)
	itemProduct.IDUpdatedby = userID

	err = h.itemProductService.Update(itemProduct)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": itemProduct.ID,
	}

	return pkg.Response(c, fiber.StatusOK, "Item Product Type updated successfully", result)
}

// DeleteItemProductType godoc
// @Summary Delete a itemProduct
// @Description Delete a Item Product Type by its ID
// @Tags Item Product Type
// @Param id path int true "Item Product Type ID"
// @Success 200 {object} map[string]interface{} "Item Product Type deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Item Product Type not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /general/master/item/product-type/{id} [delete]
func (h *ItemProductTypeHandler) DeleteItemProductType(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	user, err := h.itemProductService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Item Product Type not found"))
	}

	err = h.itemProductService.Delete(user)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Item Product Type deleted successfully", nil)
}
