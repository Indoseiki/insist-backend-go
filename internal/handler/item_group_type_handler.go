package handler

import (
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type ItemGroupTypeHandler struct {
	itemGroupService *service.ItemGroupTypeService
}

func NewItemGroupTypeHandler(itemGroupService *service.ItemGroupTypeService) *ItemGroupTypeHandler {
	return &ItemGroupTypeHandler{itemGroupService: itemGroupService}
}

// GetItemGroupTypes godoc
// @Summary Get a list of Item Group Types
// @Description Retrieves Item Group Types with pagination, search, filtering, and sorting options
// @Tags Item Group Type
// @Accept json
// @Produce json
// @Param idItemGroup query int false "Filter by Item Group ID"
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering Item Group Type"
// @Param sortBy query string false "Field to sort by"
// @Param sortDirection query boolean false "Sort direction: true for ascending, false for descending"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /general/master/item/group-type [get]
func (h *ItemGroupTypeHandler) GetItemGroupTypes(c *fiber.Ctx) error {
	idItemGroup := c.QueryInt("idItemGroup", 0)
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	total, err := h.itemGroupService.GetTotal(search, uint(idItemGroup))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	itemGroupTypes, err := h.itemGroupService.GetAll(offset, rows, search, sortBy, sortDirection, uint(idItemGroup))
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
		"items": itemGroupTypes,
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

	if len(itemGroupTypes) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetItemGroupType godoc
// @Summary Get Item Group Type by ID
// @Description Retrieve a specific Item Group Type by its ID
// @Tags Item Group Type
// @Accept json
// @Produce json
// @Param id path int true "Item Group Type ID"
// @Success 200 {object} map[string]interface{} "Item Group Type found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Item Group Type not found"
// @Router /general/master/item/group-type/{id} [get]
func (h *ItemGroupTypeHandler) GetItemGroupType(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	itemGroup, err := h.itemGroupService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Item Group Type not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Item Group Type found successfully", itemGroup)
}

// CreateItemGroupType godoc
// @Summary Create a new Item Group Type
// @Description Create a new Item Group Type with the provided details
// @Tags Item Group Type
// @Accept json
// @Produce json
// @Param ItemGroupType body model.MstItemGroupType true "Item Group Type details"
// @Success 201 {object} map[string]interface{} "Item Group Type created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /general/master/item/group-type [post]
func (h *ItemGroupTypeHandler) CreateItemGroupType(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var itemGroup model.MstItemGroupType
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

	return pkg.Response(c, fiber.StatusCreated, "Item Group Type created successfully", result)
}

// UpdateItemGroupType godoc
// @Summary Update an existing Item Group Type
// @Description Update the details of an existing Item Group Type by its ID
// @Tags Item Group Type
// @Accept json
// @Produce json
// @Param id path int true "Item Group Type ID"
// @Param ItemGroupType body model.MstItemGroupType true "Updated Item Group Type details"
// @Success 200 {object} map[string]interface{} "Item Group Type updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 404 {object} map[string]interface{} "Not Found: Item Group Type not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /general/master/item/group-type/{id} [put]
func (h *ItemGroupTypeHandler) UpdateItemGroupType(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	itemGroup, err := h.itemGroupService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Item Group Type not found"))
	}

	if err := c.BodyParser(itemGroup); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	itemGroup.ID = uint(ID)
	itemGroup.IDUpdatedby = userID

	err = h.itemGroupService.Update(itemGroup)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": itemGroup.ID,
	}

	return pkg.Response(c, fiber.StatusOK, "Item Group Type updated successfully", result)
}

// DeleteItemGroupType godoc
// @Summary Delete an Item Group Type
// @Description Delete an Item Group Type by its ID
// @Tags Item Group Type
// @Param id path int true "Item Group Type ID"
// @Success 200 {object} map[string]interface{} "Item Group Type deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Item Group Type not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /general/master/item/group-type/{id} [delete]
func (h *ItemGroupTypeHandler) DeleteItemGroupType(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	itemGroup, err := h.itemGroupService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Item Group Type not found"))
	}

	err = h.itemGroupService.Delete(itemGroup)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Item Group Type deleted successfully", nil)
}
