package handler

import (
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type ItemGroupHandler struct {
	itemGroupService *service.ItemGroupService
}

func NewItemGroupHandler(itemGroupService *service.ItemGroupService) *ItemGroupHandler {
	return &ItemGroupHandler{itemGroupService: itemGroupService}
}

// GetItemGroups godoc
// @Summary Get a list of Item Groups
// @Description Retrieves Item Groups with pagination, filtering, and sorting options
// @Tags Item Group
// @Accept json
// @Produce json
// @Param id_product_type query int false "Filter by Product Type ID"
// @Param categoryCode query string false "Filter by Item Category Code"
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering Item Group"
// @Param sortBy query string false "Field to sort by"
// @Param sortDirection query boolean false "Sort direction: true for ascending, false for descending"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /general/master/item/group [get]
func (h *ItemGroupHandler) GetItemGroups(c *fiber.Ctx) error {
	idProductType := c.QueryInt("idProductType", 0)
	categoryCode := c.Query("categoryCode", "")
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	total, err := h.itemGroupService.GetTotal(search, categoryCode, uint(idProductType))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	itemGroups, err := h.itemGroupService.GetAll(offset, rows, search, sortBy, sortDirection, categoryCode, uint(idProductType))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	totalPages := int(math.Ceil(float64(total) / float64(rows)))

	var start *int
	if total == 0 {
		start = nil
	} else {
		val := offset + 1
		start = &val
	}

	var end *int
	if total == 0 {
		end = nil
	} else {
		val := int(math.Min(float64(offset+rows), float64(total)))
		end = &val
	}

	var nextPage *int
	if page < totalPages {
		val := page + 1
		nextPage = &val
	}

	result := map[string]interface{}{
		"items": itemGroups,
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

	if len(itemGroups) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetItemGroup godoc
// @Summary Get Item Group by ID
// @Description Retrieve a specific Item Group by its ID
// @Tags Item Group
// @Accept json
// @Produce json
// @Param id path int true "Item Group ID"
// @Success 200 {object} map[string]interface{} "Item Group found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Item Group not found"
// @Router /general/master/item/group/{id} [get]
func (h *ItemGroupHandler) GetItemGroup(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	itemGroup, err := h.itemGroupService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Item Group not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Item Group found successfully", itemGroup)
}

// CreateItemGroup godoc
// @Summary Create a new Item Group
// @Description Create a new Item Group with the provided details
// @Tags Item Group
// @Accept json
// @Produce json
// @Param ItemGroup body model.MstItemGroup true "Item Group details"
// @Success 201 {object} map[string]interface{} "Item Group created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /general/master/item/group [post]
func (h *ItemGroupHandler) CreateItemGroup(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var itemGroup model.MstItemGroup
	if err := c.BodyParser(&itemGroup); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	itemGroup.IDCreatedby = userID
	itemGroup.IDUpdatedby = userID

	err := h.itemGroupService.Create(&itemGroup)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": itemGroup.ID,
	}

	return pkg.Response(c, fiber.StatusCreated, "Item Group created successfully", result)
}

// UpdateItemGroup godoc
// @Summary Update an existing Item Group
// @Description Update the details of an existing Item Group by its ID
// @Tags Item Group
// @Accept json
// @Produce json
// @Param id path int true "Item Group ID"
// @Param ItemGroup body model.MstItemGroup true "Updated Item Group details"
// @Success 200 {object} map[string]interface{} "Item Group updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 404 {object} map[string]interface{} "Not Found: Item Group not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /general/master/item/group/{id} [put]
func (h *ItemGroupHandler) UpdateItemGroup(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	itemGroup, err := h.itemGroupService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Item Group not found"))
	}

	if err := c.BodyParser(itemGroup); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	itemGroup.IDUpdatedby = userID

	err = h.itemGroupService.Update(itemGroup)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": itemGroup.ID,
	}

	return pkg.Response(c, fiber.StatusOK, "Item Group updated successfully", result)
}

// DeleteItemGroup godoc
// @Summary Delete an Item Group
// @Description Delete an Item Group by its ID
// @Tags Item Group
// @Param id path int true "Item Group ID"
// @Success 200 {object} map[string]interface{} "Item Group deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Item Group not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /general/master/item/group/{id} [delete]
func (h *ItemGroupHandler) DeleteItemGroup(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	itemGroup, err := h.itemGroupService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Item Group not found"))
	}

	err = h.itemGroupService.Delete(itemGroup)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Item Group deleted successfully", nil)
}
