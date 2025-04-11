package handler

import (
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type ItemSubCategoryHandler struct {
	itemSubCategoryService *service.ItemSubCategoryService
}

func NewItemSubCategoryHandler(itemSubCategoryService *service.ItemSubCategoryService) *ItemSubCategoryHandler {
	return &ItemSubCategoryHandler{itemSubCategoryService: itemSubCategoryService}
}

// GetItemSubCategories godoc
// @Summary Get a list of Item Sub Categories
// @Description Retrieves Item Sub Categories with pagination, filtering, search, and sorting
// @Tags Item Sub Category
// @Accept json
// @Produce json
// @Param categoryCode query string false "Filter by Item Category Code"
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering Item Sub Category"
// @Param sortBy query string false "Field to sort by"
// @Param sortDirection query boolean false "Sort direction: true for ascending, false for descending"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /general/master/item/sub-category [get]
func (h *ItemSubCategoryHandler) GetItemSubCategories(c *fiber.Ctx) error {
	categoryCode := c.Query("categoryCode", "")
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	total, err := h.itemSubCategoryService.GetTotal(search, categoryCode)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	itemSubCategories, err := h.itemSubCategoryService.GetAll(offset, rows, search, sortBy, sortDirection, categoryCode)
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
		"items": itemSubCategories,
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

	if len(itemSubCategories) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetItemSubCategory godoc
// @Summary Get Item Sub Category by ID
// @Description Retrieve a specific Item Sub Category by its ID
// @Tags Item Sub Category
// @Accept json
// @Produce json
// @Param id path int true "Item Sub Category ID"
// @Success 200 {object} map[string]interface{} "Item Sub Category found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Item Sub Category not found"
// @Router /general/master/item/sub-category/{id} [get]
func (h *ItemSubCategoryHandler) GetItemSubCategory(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	itemSubCategory, err := h.itemSubCategoryService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Item Sub Category not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Item Sub Category found successfully", itemSubCategory)
}

// CreateItemSubCategory godoc
// @Summary Create a new Item Sub Category
// @Description Create a new Item Sub Category with the provided details
// @Tags Item Sub Category
// @Accept json
// @Produce json
// @Param ItemSubCategory body model.MstItemSubCategory true "Item Sub Category details"
// @Success 201 {object} map[string]interface{} "Item Sub Category created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /general/master/item/sub-category [post]
func (h *ItemSubCategoryHandler) CreateItemSubCategory(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var itemSubCategory model.MstItemSubCategory
	if err := c.BodyParser(&itemSubCategory); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	itemSubCategory.IDCreatedby = userID
	itemSubCategory.IDUpdatedby = userID

	err := h.itemSubCategoryService.Create(&itemSubCategory)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": itemSubCategory.ID,
	}

	return pkg.Response(c, fiber.StatusCreated, "Item Sub Category created successfully", result)
}

// UpdateItemSubCategory godoc
// @Summary Update an existing Item Sub Category
// @Description Update the details of an existing Item Sub Category by its ID
// @Tags Item Sub Category
// @Accept json
// @Produce json
// @Param id path int true "Item Sub Category ID"
// @Param ItemSubCategory body model.MstItemSubCategory true "Updated Item Sub Category details"
// @Success 200 {object} map[string]interface{} "Item Sub Category updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 404 {object} map[string]interface{} "Not Found: Item Sub Category not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /general/master/item/sub-category/{id} [put]
func (h *ItemSubCategoryHandler) UpdateItemSubCategory(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var itemSubCategory *model.MstItemSubCategory
	itemSubCategory, err = h.itemSubCategoryService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Item Sub Category not found"))
	}

	if err := c.BodyParser(itemSubCategory); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	itemSubCategory.ID = uint(ID)
	itemSubCategory.IDUpdatedby = userID

	err = h.itemSubCategoryService.Update(itemSubCategory)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": itemSubCategory.ID,
	}

	return pkg.Response(c, fiber.StatusOK, "Item Sub Category updated successfully", result)
}

// DeleteItemSubCategory godoc
// @Summary Delete a itemSubCategory
// @Description Delete a Item Sub Category by its ID
// @Tags Item Sub Category
// @Param id path int true "Item Sub Category ID"
// @Success 200 {object} map[string]interface{} "Item Sub Category deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Item Sub Category not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /general/master/item/sub-category/{id} [delete]
func (h *ItemSubCategoryHandler) DeleteItemSubCategory(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	user, err := h.itemSubCategoryService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Item Sub Category not found"))
	}

	err = h.itemSubCategoryService.Delete(user)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Item Sub Category deleted successfully", nil)
}
