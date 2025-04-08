package handler

import (
	"encoding/json"
	"fmt"
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type CurrencyHandler struct {
	currencyService *service.CurrencyService
}

func NewCurrencyHandler(currencyService *service.CurrencyService) *CurrencyHandler {
	return &CurrencyHandler{currencyService: currencyService}
}

// GetCurrencies godoc
// @Summary Get a list of currencies
// @Description Retrieves currencies with pagination and optional search
// @Tags Currency
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering currency"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /acf/master/currency [get]
func (h *CurrencyHandler) GetCurrencies(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	total, err := h.currencyService.GetTotal(search)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	currencies, err := h.currencyService.GetAll(offset, rows, search, sortBy, sortDirection)
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
		"items": currencies,
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

	if len(currencies) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetCurrency godoc
// @Summary Get currency by ID
// @Description Retrieve a specific currency by its ID
// @Tags Currency
// @Accept json
// @Produce json
// @Param id path int true "Currency ID"
// @Success 200 {object} map[string]interface{} "Currency found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Currency not found"
// @Router /acf/master/currency/{id} [get]
func (h *CurrencyHandler) GetCurrency(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	currency, err := h.currencyService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Currency not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Currency found successfully", currency)
}

// CreateCurrency godoc
// @Summary Create a new currency
// @Description Create a new currency with the provided details
// @Tags Currency
// @Accept json
// @Produce json
// @Param currency body model.MstCurrency true "Currency details"
// @Success 201 {object} map[string]interface{} "Currency created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /acf/master/currency [post]
func (h *CurrencyHandler) CreateCurrency(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var currency model.MstCurrency
	if err := c.BodyParser(&currency); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	currency.IDCreatedby = userID
	currency.IDUpdatedby = userID

	err := h.currencyService.Create(&currency)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": currency.ID,
	}

	return pkg.Response(c, fiber.StatusCreated, "Currency created successfully", result)
}

// UpdateCurrency godoc
// @Summary Update an existing currency
// @Description Update the details of an existing currency by its ID
// @Tags Currency
// @Accept json
// @Produce json
// @Param id path int true "Currency ID"
// @Param currency body model.MstCurrency true "Updated currency details"
// @Success 200 {object} map[string]interface{} "Currency updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 404 {object} map[string]interface{} "Not Found: Currency not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /acf/master/currency/{id} [put]
func (h *CurrencyHandler) UpdateCurrency(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var currency *model.MstCurrency
	currency, err = h.currencyService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Currency not found"))
	}

	if err := c.BodyParser(currency); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	currency.ID = uint(ID)
	currency.IDUpdatedby = userID

	err = h.currencyService.Update(currency)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": currency.ID,
	}

	return pkg.Response(c, fiber.StatusOK, "Currency updated successfully", result)
}

// DeleteCurrency godoc
// @Summary Delete a currency
// @Description Delete a currency by its ID
// @Tags Currency
// @Param id path int true "Currency ID"
// @Success 200 {object} map[string]interface{} "Currency deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Currency not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /acf/master/currency/{id} [delete]
func (h *CurrencyHandler) DeleteCurrency(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	user, err := h.currencyService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Currency not found"))
	}

	err = h.currencyService.Delete(user)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Currency deleted successfully", nil)
}

// GenerateCurrency godoc
// @Summary Generate currencies from Frankfurter API
// @Description Fetch currency data from external API and store it in the database
// @Tags Currency
// @Param id path int true "User ID who initiates the operation"
// @Success 200 {object} map[string]interface{} "Currencies generated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid user ID"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /acf/master/currency-generate [get]
func (h *CurrencyHandler) GenerateCurrency(c *fiber.Ctx) error {
	userID := uint(1)

	resp, err := http.Get("https://api.frankfurter.app/currencies")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch currency data"))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return pkg.ErrorResponse(c, fiber.NewError(resp.StatusCode, "Failed to fetch currency data"))
	}

	var currencyMap map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&currencyMap); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, "Failed to parse currency data"))
	}

	for code, description := range currencyMap {

		existing, _ := h.currencyService.GetByCurrencyCode(code)
		if existing != nil {
			continue
		}

		newCurrency := &model.MstCurrency{
			Currency:    code,
			Description: description,
			IDCreatedby: userID,
			IDUpdatedby: userID,
		}

		if err := h.currencyService.Create(newCurrency); err != nil {
			return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, fmt.Sprintf("Failed to insert currency code %s: %v", code, err)))
		}

	}

	return pkg.Response(c, fiber.StatusOK, "Currencies generated successfully", nil)
}
