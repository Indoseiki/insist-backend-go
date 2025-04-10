package handler

import (
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type ItemCategoryHandler struct {
	itemCategoryService *service.ItemCategoryService
}

func NewItemCategoryHandler(itemCategoryService *service.ItemCategoryService) *ItemCategoryHandler {
	return &ItemCategoryHandler{itemCategoryService: itemCategoryService}
}

// GetItemCategories godoc
// @Summary Get a list of Item Categories
// @Description Retrieves Item Categories with pagination and optional search
// @Tags Item Category
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering itemCategory"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/item-category [get]
func (h *ItemCategoryHandler) GetItemCategories(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	total, err := h.itemCategoryService.GetTotal(search)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	itemCategories, err := h.itemCategoryService.GetAll(offset, rows, search, sortBy, sortDirection)
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
		"items": itemCategories,
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

	if len(itemCategories) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetItemCategory godoc
// @Summary Get Item Category by ID
// @Description Retrieve a specific Item Category by its ID
// @Tags Item Category
// @Accept json
// @Produce json
// @Param id path int true "Item Category ID"
// @Success 200 {object} map[string]interface{} "Item Category found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Item Category not found"
// @Router /admin/master/item-category/{id} [get]
func (h *ItemCategoryHandler) GetItemCategory(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	itemCategory, err := h.itemCategoryService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Item Category not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Item Category found successfully", itemCategory)
}

// CreateItemCategory godoc
// @Summary Create a new Item Category
// @Description Create a new Item Category with the provided details
// @Tags Item Category
// @Accept json
// @Produce json
// @Param ItemCategory body model.MstItemCategory true "Item Category details"
// @Success 201 {object} map[string]interface{} "Item Category created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/item-category [post]
func (h *ItemCategoryHandler) CreateItemCategory(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var itemCategory model.MstItemCategory
	if err := c.BodyParser(&itemCategory); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	itemCategory.IDCreatedby = userID
	itemCategory.IDUpdatedby = userID

	err := h.itemCategoryService.Create(&itemCategory)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": itemCategory.ID,
	}

	return pkg.Response(c, fiber.StatusCreated, "Item Category created successfully", result)
}

// UpdateItemCategory godoc
// @Summary Update an existing Item Category
// @Description Update the details of an existing Item Category by its ID
// @Tags Item Category
// @Accept json
// @Produce json
// @Param id path int true "Item Category ID"
// @Param ItemCategory body model.MstItemCategory true "Updated Item Category details"
// @Success 200 {object} map[string]interface{} "Item Category updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 404 {object} map[string]interface{} "Not Found: Item Category not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/item-category/{id} [put]
func (h *ItemCategoryHandler) UpdateItemCategory(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var itemCategory *model.MstItemCategory
	itemCategory, err = h.itemCategoryService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Item Category not found"))
	}

	if err := c.BodyParser(itemCategory); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	itemCategory.ID = uint(ID)
	itemCategory.IDUpdatedby = userID

	err = h.itemCategoryService.Update(itemCategory)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": itemCategory.ID,
	}

	return pkg.Response(c, fiber.StatusOK, "Item Category updated successfully", result)
}

// DeleteItemCategory godoc
// @Summary Delete a itemCategory
// @Description Delete a Item Category by its ID
// @Tags Item Category
// @Param id path int true "Item Category ID"
// @Success 200 {object} map[string]interface{} "Item Category deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Item Category not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/item-category/{id} [delete]
func (h *ItemCategoryHandler) DeleteItemCategory(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	user, err := h.itemCategoryService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Item Category not found"))
	}

	err = h.itemCategoryService.Delete(user)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Item Category deleted successfully", nil)
}
