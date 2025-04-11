package handler

import (
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type ItemProductHandler struct {
	itemProductService *service.ItemProductService
}

func NewItemProductHandler(itemProductService *service.ItemProductService) *ItemProductHandler {
	return &ItemProductHandler{itemProductService: itemProductService}
}

// GetItemProducts godoc
// @Summary Get a list of Item Products
// @Description Retrieves Item Products with pagination, search, filtering, and sorting options
// @Tags Item Product
// @Accept json
// @Produce json
// @Param categoryCode query string false "Filter by Item Category Code"
// @Param subCategoryCode query string false "Filter by Item Sub Category Code"
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering Item Product"
// @Param sortBy query string false "Field to sort by"
// @Param sortDirection query boolean false "Sort direction: true for ascending, false for descending"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /general/master/item/product [get]
func (h *ItemProductHandler) GetItemProducts(c *fiber.Ctx) error {
	categoryCode := c.Query("categoryCode", "")
	subCategoryCode := c.Query("subCategoryCode", "")
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	total, err := h.itemProductService.GetTotal(search, categoryCode, subCategoryCode)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	itemProducts, err := h.itemProductService.GetAll(offset, rows, search, sortBy, sortDirection, categoryCode, subCategoryCode)
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
		"items": itemProducts,
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

	if len(itemProducts) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetItemProduct godoc
// @Summary Get Item Product by ID
// @Description Retrieve a specific Item Product by its ID
// @Tags Item Product
// @Accept json
// @Produce json
// @Param id path int true "Item Product ID"
// @Success 200 {object} map[string]interface{} "Item Product found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Item Product not found"
// @Router /general/master/item/product/{id} [get]
func (h *ItemProductHandler) GetItemProduct(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	itemProduct, err := h.itemProductService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Item Product not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Item Product found successfully", itemProduct)
}

// CreateItemProduct godoc
// @Summary Create a new Item Product
// @Description Create a new Item Product with the provided details
// @Tags Item Product
// @Accept json
// @Produce json
// @Param ItemProduct body model.MstItemProduct true "Item Product details"
// @Success 201 {object} map[string]interface{} "Item Product created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /general/master/item/product [post]
func (h *ItemProductHandler) CreateItemProduct(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var itemProduct model.MstItemProduct
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

	return pkg.Response(c, fiber.StatusCreated, "Item Product created successfully", result)
}

// UpdateItemProduct godoc
// @Summary Update an existing Item Product
// @Description Update the details of an existing Item Product by its ID
// @Tags Item Product
// @Accept json
// @Produce json
// @Param id path int true "Item Product ID"
// @Param ItemProduct body model.MstItemProduct true "Updated Item Product details"
// @Success 200 {object} map[string]interface{} "Item Product updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 404 {object} map[string]interface{} "Not Found: Item Product not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /general/master/item/product/{id} [put]
func (h *ItemProductHandler) UpdateItemProduct(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var itemProduct *model.MstItemProduct
	itemProduct, err = h.itemProductService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Item Product not found"))
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

	return pkg.Response(c, fiber.StatusOK, "Item Product updated successfully", result)
}

// DeleteItemProduct godoc
// @Summary Delete a itemProduct
// @Description Delete a Item Product by its ID
// @Tags Item Product
// @Param id path int true "Item Product ID"
// @Success 200 {object} map[string]interface{} "Item Product deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Item Product not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /general/master/item/product/{id} [delete]
func (h *ItemProductHandler) DeleteItemProduct(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	user, err := h.itemProductService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Item Product not found"))
	}

	err = h.itemProductService.Delete(user)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Item Product deleted successfully", nil)
}
