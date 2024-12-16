package handler

import (
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type KeyValueHandler struct {
	keyValueService *service.KeyValueService
}

func NewKeyValueHandler(keyValueService *service.KeyValueService) *KeyValueHandler {
	return &KeyValueHandler{keyValueService: keyValueService}
}

// GetKeyValues godoc
// @Summary Get a list of key values
// @Description Retrieves key values with pagination and optional search
// @Tags Key Value
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering key value"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/key-value [get]
func (h *KeyValueHandler) GetKeyValues(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search")
	offset := (page - 1) * rows

	total, err := h.keyValueService.GetTotal(search)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	keyValues, err := h.keyValueService.GetAll(offset, rows, search)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	totalPages := int(math.Ceil(float64(total) / float64(rows)))
	var nextPage *int
	if page < totalPages {
		nextPageVal := page + 1
		nextPage = &nextPageVal
	}

	result := map[string]interface{}{
		"items": keyValues,
		"pagination": map[string]interface{}{
			"current_page":  page,
			"next_page":     nextPage,
			"total_pages":   totalPages,
			"rows_per_page": rows,
			"total_rows":    total,
		},
	}

	if len(keyValues) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetKeyValue godoc
// @Summary Get key value by ID
// @Description Retrieve a specific key value by its ID
// @Tags Key Value
// @Accept json
// @Produce json
// @Param id path int true "Key Value ID"
// @Success 200 {object} map[string]interface{} "Key Value found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Key Value not found"
// @Router /admin/master/key-value/{id} [get]
func (h *KeyValueHandler) GetKeyValue(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	keyValue, err := h.keyValueService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Key Value not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Key Value found successfully", keyValue)
}

// CreateKeyValue godoc
// @Summary Create a new key value
// @Description Create a new key value with the provided details
// @Tags Key Value
// @Accept json
// @Produce json
// @Param keyValue body model.MstKeyValue true "Key Value details"
// @Success 201 {object} map[string]interface{} "Key Value created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/key-value [post]
func (h *KeyValueHandler) CreateKeyValue(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var keyValue model.MstKeyValue
	if err := c.BodyParser(&keyValue); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	keyValue.IDCreatedby = userID
	keyValue.IDUpdatedby = userID

	err := h.keyValueService.Create(&keyValue)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": keyValue.ID,
	}

	return pkg.Response(c, fiber.StatusCreated, "Key Value created successfully", result)
}

// UpdateKeyValue godoc
// @Summary Update an existing key value
// @Description Update the details of an existing key value by its ID
// @Tags Key Value
// @Accept json
// @Produce json
// @Param id path int true "Key Value ID"
// @Param keyValue body model.MstKeyValue true "Updated key value details"
// @Success 200 {object} map[string]interface{} "Key Value updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 404 {object} map[string]interface{} "Not Found: Key Value not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/key-value/{id} [put]
func (h *KeyValueHandler) UpdateKeyValue(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var keyValue *model.MstKeyValue
	keyValue, err = h.keyValueService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Key Value not found"))
	}

	if err := c.BodyParser(keyValue); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	keyValue.ID = uint(ID)
	keyValue.IDUpdatedby = userID

	err = h.keyValueService.Update(keyValue)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": keyValue.ID,
	}

	return pkg.Response(c, fiber.StatusOK, "Key Value updated successfully", result)
}

// DeleteKeyValue godoc
// @Summary Delete a key value
// @Description Delete a key value by its ID
// @Tags Key Value
// @Param id path int true "Key Value ID"
// @Success 200 {object} map[string]interface{} "Key Value deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Key Value not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/key-value/{id} [delete]
func (h *KeyValueHandler) DeleteKeyValue(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	user, err := h.keyValueService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Key Value not found"))
	}

	err = h.keyValueService.Delete(user)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Key Value deleted successfully", nil)
}
