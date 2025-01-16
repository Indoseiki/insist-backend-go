package handler

import (
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type WarehouseHandler struct {
	warehouseService *service.WarehouseService
}

func NewWarehouseHandler(warehouseService *service.WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{warehouseService: warehouseService}
}

// GetWarehouses godoc
// @Summary Get a list of warehouses
// @Description Retrieves warehouses with pagination and optional search
// @Tags Warehouse
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering warehouse"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /pid/master/warehouse [get]
func (h *WarehouseHandler) GetWarehouses(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	total, err := h.warehouseService.GetTotal(search)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	warehouses, err := h.warehouseService.GetAll(offset, rows, search, sortBy, sortDirection)
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
		"items": warehouses,
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

	if len(warehouses) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetWarehouse godoc
// @Summary Get warehouse by ID
// @Description Retrieve a specific warehouse by its ID
// @Tags Warehouse
// @Accept json
// @Produce json
// @Param id path int true "Warehouse ID"
// @Success 200 {object} map[string]interface{} "Warehouse found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Warehouse not found"
// @Router /pid/master/warehouse/{id} [get]
func (h *WarehouseHandler) GetWarehouse(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	warehouse, err := h.warehouseService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Warehouse not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Warehouse found successfully", warehouse)
}

// CreateWarehouse godoc
// @Summary Create a new warehouse
// @Description Create a new warehouse with the provided details
// @Tags Warehouse
// @Accept json
// @Produce json
// @Param warehouse body model.MstWarehouse true "Warehouse details"
// @Success 201 {object} map[string]interface{} "Warehouse created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /pid/master/warehouse [post]
func (h *WarehouseHandler) CreateWarehouse(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var warehouse model.MstWarehouse
	if err := c.BodyParser(&warehouse); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	warehouse.IDCreatedby = userID
	warehouse.IDUpdatedby = userID

	err := h.warehouseService.Create(&warehouse)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": warehouse.ID,
	}

	return pkg.Response(c, fiber.StatusCreated, "Warehouse created successfully", result)
}

// UpdateWarehouse godoc
// @Summary Update an existing warehouse
// @Description Update the details of an existing warehouse by its ID
// @Tags Warehouse
// @Accept json
// @Produce json
// @Param id path int true "Warehouse ID"
// @Param warehouse body model.MstWarehouse true "Updated warehouse details"
// @Success 200 {object} map[string]interface{} "Warehouse updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 404 {object} map[string]interface{} "Not Found: Warehouse not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /pid/master/warehouse/{id} [put]
func (h *WarehouseHandler) UpdateWarehouse(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var warehouse *model.MstWarehouse
	warehouse, err = h.warehouseService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Warehouse not found"))
	}

	if err := c.BodyParser(warehouse); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	warehouse.ID = uint(ID)
	warehouse.IDUpdatedby = userID

	err = h.warehouseService.Update(warehouse)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": warehouse.ID,
	}

	return pkg.Response(c, fiber.StatusOK, "Warehouse updated successfully", result)
}

// DeleteWarehouse godoc
// @Summary Delete a warehouse
// @Description Delete a warehouse by its ID
// @Tags Warehouse
// @Param id path int true "Warehouse ID"
// @Success 200 {object} map[string]interface{} "Warehouse deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Warehouse not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /pid/master/warehouse/{id} [delete]
func (h *WarehouseHandler) DeleteWarehouse(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	user, err := h.warehouseService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Warehouse not found"))
	}

	err = h.warehouseService.Delete(user)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Warehouse deleted successfully", nil)
}
