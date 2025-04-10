package handler

import (
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type BillingTermHandler struct {
	billingTermService *service.BillingTermService
}

func NewBillingTermHandler(billingTermService *service.BillingTermService) *BillingTermHandler {
	return &BillingTermHandler{billingTermService: billingTermService}
}

// GetBillingTerms godoc
// @Summary Get a list of billing terms
// @Description Retrieves billing terms with pagination and optional search
// @Tags BillingTerm
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering billing term"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /general/master/billing-term [get]
func (h *BillingTermHandler) GetBillingTerms(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	total, err := h.billingTermService.GetTotal(search)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	billingTerms, err := h.billingTermService.GetAll(offset, rows, search, sortBy, sortDirection)
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
		"items": billingTerms,
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

	if len(billingTerms) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetBillingTerm godoc
// @Summary Get billing term by ID
// @Description Retrieve a specific billing term by its ID
// @Tags BillingTerm
// @Accept json
// @Produce json
// @Param id path int true "Billing Term ID"
// @Success 200 {object} map[string]interface{} "Billing Term found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Billing Term not found"
// @Router /general/master/billing-term/{id} [get]
func (h *BillingTermHandler) GetBillingTerm(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	billingTerm, err := h.billingTermService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Billing Term not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Billing Term found successfully", billingTerm)
}

// CreateBillingTerm godoc
// @Summary Create a new billing term
// @Description Create a new billing term with the provided details
// @Tags BillingTerm
// @Accept json
// @Produce json
// @Param billingTerm body model.MstBillingTerm true "Billing Term details"
// @Success 201 {object} map[string]interface{} "Billing Term created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /general/master/billing-term [post]
func (h *BillingTermHandler) CreateBillingTerm(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var billingTerm model.MstBillingTerm
	if err := c.BodyParser(&billingTerm); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	billingTerm.IDCreatedby = userID
	billingTerm.IDUpdatedby = userID

	err := h.billingTermService.Create(&billingTerm)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": billingTerm.ID,
	}

	return pkg.Response(c, fiber.StatusCreated, "Billing Term created successfully", result)
}

// UpdateBillingTerm godoc
// @Summary Update an existing billing term
// @Description Update the details of an existing billing term by its ID
// @Tags BillingTerm
// @Accept json
// @Produce json
// @Param id path int true "Billing Term ID"
// @Param billingTerm body model.MstBillingTerm true "Updated billing term details"
// @Success 200 {object} map[string]interface{} "Billing Term updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 404 {object} map[string]interface{} "Not Found: Billing Term not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /general/master/billing-term/{id} [put]
func (h *BillingTermHandler) UpdateBillingTerm(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var billingTerm *model.MstBillingTerm
	billingTerm, err = h.billingTermService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Billing Term not found"))
	}

	if err := c.BodyParser(billingTerm); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	billingTerm.ID = uint(ID)
	billingTerm.IDUpdatedby = userID

	err = h.billingTermService.Update(billingTerm)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": billingTerm.ID,
	}

	return pkg.Response(c, fiber.StatusOK, "Billing Term updated successfully", result)
}

// DeleteBillingTerm godoc
// @Summary Delete a billing term
// @Description Delete a billing term by its ID
// @Tags BillingTerm
// @Param id path int true "Billing Term ID"
// @Success 200 {object} map[string]interface{} "Billing Term deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Billing Term not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /general/master/billing-term/{id} [delete]
func (h *BillingTermHandler) DeleteBillingTerm(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	user, err := h.billingTermService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Billing Term not found"))
	}

	err = h.billingTermService.Delete(user)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Billing Term deleted successfully", nil)
}
