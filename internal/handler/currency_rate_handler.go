package handler

import (
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type CurrencyRateHandler struct {
	currencyService *service.CurrencyRateService
}

func NewCurrencyRateHandler(currencyService *service.CurrencyRateService) *CurrencyRateHandler {
	return &CurrencyRateHandler{currencyService: currencyService}
}

// GetCurrencyRates godoc
// @Summary Get a list of currency rates
// @Description Retrieves currency rates with pagination and optional search
// @Tags Currency Rate
// @Accept json
// @Produce json
// @Param idCurrency query int false "Currency ID"
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering currency"
// @Param sortBy query string false "Field to sort by"
// @Param sortDirection query bool false "Sorting direction (true for ascending, false for descending)"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /acf/master/currency-rate [get]
func (h *CurrencyRateHandler) GetCurrencyRates(c *fiber.Ctx) error {
	idCurrency := c.QueryInt("idCurrency", 0)
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	total, err := h.currencyService.GetTotal(uint(idCurrency), search)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	currencyRates, err := h.currencyService.GetAll(uint(idCurrency), offset, rows, search, sortBy, sortDirection)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	totalPages := int(math.Ceil(float64(total) / float64(rows)))

	var start, end, nextPage *int
	if total > 0 {
		startVal := offset + 1
		start = &startVal
		endVal := int(math.Min(float64(offset+rows), float64(total)))
		end = &endVal
		if page < totalPages {
			nextPageVal := page + 1
			nextPage = &nextPageVal
		}
	}

	result := map[string]interface{}{
		"items": currencyRates,
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

	if len(currencyRates) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetCurrencyRate godoc
// @Summary Get currency rate by ID
// @Description Retrieve a specific currency rate by its ID
// @Tags Currency Rate
// @Accept json
// @Produce json
// @Param id path int true "Currency Rate ID"
// @Success 200 {object} map[string]interface{} "Currency Rate found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Currency Rate not found"
// @Router /acf/master/currency-rate/{id} [get]
func (h *CurrencyRateHandler) GetCurrencyRate(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	currency, err := h.currencyService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Currency Rate not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Currency Rate found successfully", currency)
}

// CreateCurrencyRate godoc
// @Summary Create a new currency rate
// @Description Create a new currency rate with the provided details
// @Tags Currency Rate
// @Accept json
// @Produce json
// @Param currency_rate body model.MstCurrencyRate true "Currency Rate details"
// @Success 201 {object} map[string]interface{} "Currency Rate created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /acf/master/currency-rate [post]
func (h *CurrencyRateHandler) CreateCurrencyRate(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var currencyRate model.MstCurrencyRate
	if err := c.BodyParser(&currencyRate); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	currencyRate.IDCreatedby = userID
	currencyRate.IDUpdatedby = userID

	err := h.currencyService.Create(&currencyRate)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": currencyRate.ID,
	}

	return pkg.Response(c, fiber.StatusCreated, "Currency Rate created successfully", result)
}

// UpdateCurrencyRate godoc
// @Summary Update an existing currency
// @Description Update the details of an existing currency rate by its ID
// @Tags Currency Rate
// @Accept json
// @Produce json
// @Param id path int true "Currency Rate ID"
// @Param currency_rate body model.MstCurrencyRate true "Updated currency rate details"
// @Success 200 {object} map[string]interface{} "Currency Rate updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 404 {object} map[string]interface{} "Not Found: Currency Rate not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /acf/master/currency-rate/{id} [put]
func (h *CurrencyRateHandler) UpdateCurrencyRate(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var currencyRate *model.MstCurrencyRate
	currencyRate, err = h.currencyService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Currency Rate not found"))
	}

	if err := c.BodyParser(currencyRate); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	currencyRate.ID = uint(ID)
	currencyRate.IDUpdatedby = userID

	err = h.currencyService.Update(currencyRate)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": currencyRate.ID,
	}

	return pkg.Response(c, fiber.StatusOK, "Currency Rate updated successfully", result)
}

// DeleteCurrencyRate godoc
// @Summary Delete a currency
// @Description Delete a currency rate by its ID
// @Tags Currency Rate
// @Param id path int true "Currency Rate ID"
// @Success 200 {object} map[string]interface{} "Currency Rate deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Currency Rate not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /acf/master/currency-rate/{id} [delete]
func (h *CurrencyRateHandler) DeleteCurrencyRate(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	user, err := h.currencyService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Currency Rate not found"))
	}

	err = h.currencyService.Delete(user)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Currency Rate deleted successfully", nil)
}
