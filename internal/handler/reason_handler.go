package handler

import (
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type ReasonHandler struct {
	reasonService *service.ReasonService
}

func NewReasonHandler(reasonService *service.ReasonService) *ReasonHandler {
	return &ReasonHandler{reasonService: reasonService}
}

// GetReasons godoc
// @Summary Get a list of reasons
// @Description Retrieves reasons with pagination and optional search
// @Tags Reason
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param id_menu query int false "ID Menu" default(0)
// @Param search query string false "Search keyword for filtering reason"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/reason [get]
func (h *ReasonHandler) GetReasons(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	menuID := c.QueryInt("id_menu", 0)
	search := c.Query("search")
	offset := (page - 1) * rows

	total, err := h.reasonService.GetTotal(search, uint(menuID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	reasons, err := h.reasonService.GetAll(offset, rows, search, uint(menuID))
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
		"items": reasons,
		"pagination": map[string]interface{}{
			"current_page":  page,
			"next_page":     nextPage,
			"total_pages":   totalPages,
			"rows_per_page": rows,
			"total_rows":    total,
		},
	}

	if len(reasons) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetReason godoc
// @Summary Get reason by ID
// @Description Retrieve a specific reason by its ID
// @Tags Reason
// @Accept json
// @Produce json
// @Param id path int true "Reason ID"
// @Success 200 {object} map[string]interface{} "Reason found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Reason not found"
// @Router /admin/master/reason/{id} [get]
func (h *ReasonHandler) GetReason(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	reason, err := h.reasonService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Reason not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Reason found successfully", reason)
}

// CreateReason godoc
// @Summary Create a new reason
// @Description Create a new reason with the provided details
// @Tags Reason
// @Accept json
// @Produce json
// @Param reason body model.MstReason true "Reason details"
// @Success 201 {object} map[string]interface{} "Reason created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/reason [post]
func (h *ReasonHandler) CreateReason(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var reason model.MstReason
	if err := c.BodyParser(&reason); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	reason.IDCreatedby = userID
	reason.IDUpdatedby = userID

	err := h.reasonService.Create(&reason)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": reason.ID,
	}

	return pkg.Response(c, fiber.StatusCreated, "Reason created successfully", result)
}

// UpdateReason godoc
// @Summary Update an existing reason
// @Description Update the details of an existing reason by its ID
// @Tags Reason
// @Accept json
// @Produce json
// @Param id path int true "Reason ID"
// @Param reason body model.MstReason true "Updated reason details"
// @Success 200 {object} map[string]interface{} "Reason updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 404 {object} map[string]interface{} "Not Found: Reason not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/reason/{id} [put]
func (h *ReasonHandler) UpdateReason(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var reason *model.MstReason
	reason, err = h.reasonService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Reason not found"))
	}

	if err := c.BodyParser(reason); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	reason.ID = uint(ID)
	reason.IDUpdatedby = userID

	err = h.reasonService.Update(reason)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": reason.ID,
	}

	return pkg.Response(c, fiber.StatusOK, "Reason updated successfully", result)
}

// DeleteReason godoc
// @Summary Delete a reason
// @Description Delete a reason by its ID
// @Tags Reason
// @Param id path int true "Reason ID"
// @Success 200 {object} map[string]interface{} "Reason deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Reason not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/reason/{id} [delete]
func (h *ReasonHandler) DeleteReason(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	user, err := h.reasonService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Reason not found"))
	}

	err = h.reasonService.Delete(user)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Reason deleted successfully", nil)
}
