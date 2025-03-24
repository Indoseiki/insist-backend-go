package handler

import (
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type TaxCodeHandler struct {
	taxCodeService *service.TaxCodeService
}

func NewTaxCodeHandler(taxCodeService *service.TaxCodeService) *TaxCodeHandler {
	return &TaxCodeHandler{taxCodeService: taxCodeService}
}

// GetTaxCodes godoc
// @Summary Get a list of tax codes
// @Description Retrieves tax codes with pagination and optional search
// @Tags Tax Code
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering tax code"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /acf/master/tax-code [get]
func (h *TaxCodeHandler) GetTaxCodes(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	total, err := h.taxCodeService.GetTotal(search)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	taxCodes, err := h.taxCodeService.GetAll(offset, rows, search, sortBy, sortDirection)
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
		"items": taxCodes,
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

	if len(taxCodes) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetTaxCode godoc
// @Summary Get tax code by ID
// @Description Retrieve a specific tax code by its ID
// @Tags Tax Code
// @Accept json
// @Produce json
// @Param id path int true "Tax Code ID"
// @Success 200 {object} map[string]interface{} "Tax Code found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Tax Code not found"
// @Router /acf/master/tax-code/{id} [get]
func (h *TaxCodeHandler) GetTaxCode(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	taxCode, err := h.taxCodeService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Tax Code not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Tax Code found successfully", taxCode)
}

// CreateTaxCode godoc
// @Summary Create a new tax code
// @Description Create a new tax code with the provided details
// @Tags Tax Code
// @Accept json
// @Produce json
// @Param taxCode body model.MstTaxCode true "Tax Code details"
// @Success 201 {object} map[string]interface{} "Tax Code created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /acf/master/tax-code [post]
func (h *TaxCodeHandler) CreateTaxCode(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var taxCode model.MstTaxCode
	if err := c.BodyParser(&taxCode); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	taxCode.IDCreatedby = userID
	taxCode.IDUpdatedby = userID

	err := h.taxCodeService.Create(&taxCode)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": taxCode.ID,
	}

	return pkg.Response(c, fiber.StatusCreated, "Tax Code created successfully", result)
}

// UpdateTaxCode godoc
// @Summary Update an existing tax code
// @Description Update the details of an existing tax code by its ID
// @Tags Tax Code
// @Accept json
// @Produce json
// @Param id path int true "Tax Code ID"
// @Param taxCode body model.MstTaxCode true "Updated tax code details"
// @Success 200 {object} map[string]interface{} "Tax Code updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 404 {object} map[string]interface{} "Not Found: Tax Code not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /acf/master/tax-code/{id} [put]
func (h *TaxCodeHandler) UpdateTaxCode(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var taxCode *model.MstTaxCode
	taxCode, err = h.taxCodeService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Tax Code not found"))
	}

	if err := c.BodyParser(taxCode); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	taxCode.ID = uint(ID)
	taxCode.IDUpdatedby = userID

	err = h.taxCodeService.Update(taxCode)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": taxCode.ID,
	}

	return pkg.Response(c, fiber.StatusOK, "Tax Code updated successfully", result)
}

// DeleteTaxCode godoc
// @Summary Delete a tax code
// @Description Delete a tax code by its ID
// @Tags Tax Code
// @Param id path int true "Tax Code ID"
// @Success 200 {object} map[string]interface{} "Tax Code deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Tax Code not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /acf/master/tax-code/{id} [delete]
func (h *TaxCodeHandler) DeleteTaxCode(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	user, err := h.taxCodeService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Tax Code not found"))
	}

	err = h.taxCodeService.Delete(user)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Tax Code deleted successfully", nil)
}
