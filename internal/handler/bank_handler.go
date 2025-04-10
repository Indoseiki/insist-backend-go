package handler

import (
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type BankHandler struct {
	bankService *service.BankService
}

func NewBankHandler(bankService *service.BankService) *BankHandler {
	return &BankHandler{bankService: bankService}
}

// GetBanks godoc
// @Summary Get a list of banks
// @Description Retrieves banks with pagination and optional search
// @Tags Bank
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering bank"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /acf/master/bank [get]
func (h *BankHandler) GetBanks(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	total, err := h.bankService.GetTotal(search)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	banks, err := h.bankService.GetAll(offset, rows, search, sortBy, sortDirection)
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
		"items": banks,
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

	if len(banks) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetBank godoc
// @Summary Get bank by ID
// @Description Retrieve a specific bank by its ID
// @Tags Bank
// @Accept json
// @Produce json
// @Param id path int true "Bank ID"
// @Success 200 {object} map[string]interface{} "Bank found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Bank not found"
// @Router /acf/master/bank/{id} [get]
func (h *BankHandler) GetBank(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	bank, err := h.bankService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Bank not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Bank found successfully", bank)
}

// CreateBank godoc
// @Summary Create a new bank
// @Description Create a new bank with the provided details
// @Tags Bank
// @Accept json
// @Produce json
// @Param bank body model.MstBank true "Bank details"
// @Success 201 {object} map[string]interface{} "Bank created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /acf/master/bank [post]
func (h *BankHandler) CreateBank(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var bank model.MstBank
	if err := c.BodyParser(&bank); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	bank.IDCreatedby = userID
	bank.IDUpdatedby = userID

	err := h.bankService.Create(&bank)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": bank.ID,
	}

	return pkg.Response(c, fiber.StatusCreated, "Bank created successfully", result)
}

// UpdateBank godoc
// @Summary Update an existing bank
// @Description Update the details of an existing bank by its ID
// @Tags Bank
// @Accept json
// @Produce json
// @Param id path int true "Bank ID"
// @Param bank body model.MstBank true "Updated bank details"
// @Success 200 {object} map[string]interface{} "Bank updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 404 {object} map[string]interface{} "Not Found: Bank not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /acf/master/bank/{id} [put]
func (h *BankHandler) UpdateBank(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var bank *model.MstBank
	bank, err = h.bankService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Bank not found"))
	}

	if err := c.BodyParser(bank); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	bank.ID = uint(ID)
	bank.IDUpdatedby = userID

	err = h.bankService.Update(bank)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": bank.ID,
	}

	return pkg.Response(c, fiber.StatusOK, "Bank updated successfully", result)
}

// DeleteBank godoc
// @Summary Delete a bank
// @Description Delete a bank by its ID
// @Tags Bank
// @Param id path int true "Bank ID"
// @Success 200 {object} map[string]interface{} "Bank deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Bank not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /acf/master/bank/{id} [delete]
func (h *BankHandler) DeleteBank(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	user, err := h.bankService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Bank not found"))
	}

	err = h.bankService.Delete(user)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Bank deleted successfully", nil)
}
