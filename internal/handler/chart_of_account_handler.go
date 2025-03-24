package handler

import (
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type ChartOfAccountHandler struct {
	chartOfAccountService *service.ChartOfAccountService
}

func NewChartOfAccountHandler(chartOfAccountService *service.ChartOfAccountService) *ChartOfAccountHandler {
	return &ChartOfAccountHandler{chartOfAccountService: chartOfAccountService}
}

// GetChartOfAccounts godoc
// @Summary Get a list of chart of accounts
// @Description Retrieves chart of accounts with pagination and optional search
// @Tags Chart Of Account
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering chart of account"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /acf/master/chart-of-account [get]
func (h *ChartOfAccountHandler) GetChartOfAccounts(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	total, err := h.chartOfAccountService.GetTotal(search)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	chartOfAccounts, err := h.chartOfAccountService.GetAll(offset, rows, search, sortBy, sortDirection)
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
		"items": chartOfAccounts,
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

	if len(chartOfAccounts) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetChartOfAccount godoc
// @Summary Get chart of account by ID
// @Description Retrieve a specific chart of account by its ID
// @Tags Chart Of Account
// @Accept json
// @Produce json
// @Param id path int true "Chart Of Account ID"
// @Success 200 {object} map[string]interface{} "Chart Of Account found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Chart Of Account not found"
// @Router /acf/master/chart-of-account/{id} [get]
func (h *ChartOfAccountHandler) GetChartOfAccount(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	chartOfAccount, err := h.chartOfAccountService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Chart Of Account not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Chart Of Account found successfully", chartOfAccount)
}

// CreateChartOfAccount godoc
// @Summary Create a new chart of account
// @Description Create a new chart of account with the provided details
// @Tags Chart Of Account
// @Accept json
// @Produce json
// @Param chartOfAccount body model.MstChartOfAccount true "Chart Of Account details"
// @Success 201 {object} map[string]interface{} "Chart Of Account created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /acf/master/chart-of-account [post]
func (h *ChartOfAccountHandler) CreateChartOfAccount(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var chartOfAccount model.MstChartOfAccount
	if err := c.BodyParser(&chartOfAccount); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	chartOfAccount.IDCreatedby = userID
	chartOfAccount.IDUpdatedby = userID

	err := h.chartOfAccountService.Create(&chartOfAccount)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": chartOfAccount.ID,
	}

	return pkg.Response(c, fiber.StatusCreated, "Chart Of Account created successfully", result)
}

// UpdateChartOfAccount godoc
// @Summary Update an existing chart of account
// @Description Update the details of an existing chart of account by its ID
// @Tags Chart Of Account
// @Accept json
// @Produce json
// @Param id path int true "Chart Of Account ID"
// @Param chartOfAccount body model.MstChartOfAccount true "Updated chart of account details"
// @Success 200 {object} map[string]interface{} "Chart Of Account updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 404 {object} map[string]interface{} "Not Found: Chart Of Account not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /acf/master/chart-of-account/{id} [put]
func (h *ChartOfAccountHandler) UpdateChartOfAccount(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var chartOfAccount *model.MstChartOfAccount
	chartOfAccount, err = h.chartOfAccountService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Chart Of Account not found"))
	}

	if err := c.BodyParser(chartOfAccount); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	chartOfAccount.ID = uint(ID)
	chartOfAccount.IDUpdatedby = userID

	err = h.chartOfAccountService.Update(chartOfAccount)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": chartOfAccount.ID,
	}

	return pkg.Response(c, fiber.StatusOK, "Chart Of Account updated successfully", result)
}

// DeleteChartOfAccount godoc
// @Summary Delete a chart of account
// @Description Delete a chart of account by its ID
// @Tags Chart Of Account
// @Param id path int true "Chart Of Account ID"
// @Success 200 {object} map[string]interface{} "Chart Of Account deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Chart Of Account not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /acf/master/chart-of-account/{id} [delete]
func (h *ChartOfAccountHandler) DeleteChartOfAccount(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	user, err := h.chartOfAccountService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Chart Of Account not found"))
	}

	err = h.chartOfAccountService.Delete(user)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Chart Of Account deleted successfully", nil)
}
