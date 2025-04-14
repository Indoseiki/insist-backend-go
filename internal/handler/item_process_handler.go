package handler

import (
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type ItemProcessHandler struct {
	itemProcessService *service.ItemProcessService
}

func NewItemProcessHandler(itemProcessService *service.ItemProcessService) *ItemProcessHandler {
	return &ItemProcessHandler{itemProcessService: itemProcessService}
}

// GetItemProcesses godoc
// @Summary Get a list of Item Processes
// @Description Retrieves Item Processes with pagination, search, filtering, and sorting options
// @Tags Item Process
// @Accept json
// @Produce json
// @Param categoryCode query string false "Filter by Item Category Code"
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering Item Process"
// @Param sortBy query string false "Field to sort by"
// @Param sortDirection query boolean false "Sort direction: true for ascending, false for descending"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /general/master/item/process [get]
func (h *ItemProcessHandler) GetItemProcesses(c *fiber.Ctx) error {
	categoryCode := c.Query("categoryCode", "")
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	total, err := h.itemProcessService.GetTotal(search, categoryCode)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	itemProcesses, err := h.itemProcessService.GetAll(offset, rows, search, sortBy, sortDirection, categoryCode)
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
		"items": itemProcesses,
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

	if len(itemProcesses) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetItemProcess godoc
// @Summary Get Item Process by ID
// @Description Retrieve a specific Item Process by its ID
// @Tags Item Process
// @Accept json
// @Produce json
// @Param id path int true "Item Process ID"
// @Success 200 {object} map[string]interface{} "Item Process found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Item Process not found"
// @Router /general/master/item/process/{id} [get]
func (h *ItemProcessHandler) GetItemProcess(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	itemProcess, err := h.itemProcessService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Item Process not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Item Process found successfully", itemProcess)
}

// CreateItemProcess godoc
// @Summary Create a new Item Process
// @Description Create a new Item Process with the provided details
// @Tags Item Process
// @Accept json
// @Produce json
// @Param ItemProcess body model.MstItemProcess true "Item Process details"
// @Success 201 {object} map[string]interface{} "Item Process created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /general/master/item/process [post]
func (h *ItemProcessHandler) CreateItemProcess(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var itemProcess model.MstItemProcess
	if err := c.BodyParser(&itemProcess); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	itemProcess.IDCreatedby = userID
	itemProcess.IDUpdatedby = userID

	err := h.itemProcessService.Create(&itemProcess)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": itemProcess.ID,
	}

	return pkg.Response(c, fiber.StatusCreated, "Item Process created successfully", result)
}

// UpdateItemProcess godoc
// @Summary Update an existing Item Process
// @Description Update the details of an existing Item Process by its ID
// @Tags Item Process
// @Accept json
// @Produce json
// @Param id path int true "Item Process ID"
// @Param ItemProcess body model.MstItemProcess true "Updated Item Process details"
// @Success 200 {object} map[string]interface{} "Item Process updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 404 {object} map[string]interface{} "Not Found: Item Process not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /general/master/item/process/{id} [put]
func (h *ItemProcessHandler) UpdateItemProcess(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	itemProcess, err := h.itemProcessService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Item Process not found"))
	}

	if err := c.BodyParser(itemProcess); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	itemProcess.ID = uint(ID)
	itemProcess.IDUpdatedby = userID

	err = h.itemProcessService.Update(itemProcess)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": itemProcess.ID,
	}

	return pkg.Response(c, fiber.StatusOK, "Item Process updated successfully", result)
}

// DeleteItemProcess godoc
// @Summary Delete an Item Process
// @Description Delete an Item Process by its ID
// @Tags Item Process
// @Accept json
// @Produce json
// @Param id path int true "Item Process ID"
// @Success 200 {object} map[string]interface{} "Item Process deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Item Process not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /general/master/item/process/{id} [delete]
func (h *ItemProcessHandler) DeleteItemProcess(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	itemProcess, err := h.itemProcessService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Item Process not found"))
	}

	err = h.itemProcessService.Delete(itemProcess)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Item Process deleted successfully", nil)
}
