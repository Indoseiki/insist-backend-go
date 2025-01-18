package handler

import (
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type ApprovalHandler struct {
	approvalService *service.ApprovalService
}

func NewApprovalHandler(approvalService *service.ApprovalService) *ApprovalHandler {
	return &ApprovalHandler{approvalService: approvalService}
}

// GetApprovals godoc
// @Summary Get a list of approvals
// @Description Retrieves approvals with pagination and optional search
// @Tags Approval
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering approval"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/approval [get]
func (h *ApprovalHandler) GetApprovals(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search")
	offset := (page - 1) * rows

	total, err := h.approvalService.GetTotal(search)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	approvals, err := h.approvalService.GetAll(offset, rows, search)
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
		"items": approvals,
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

	if len(approvals) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetApproval godoc
// @Summary Get approval by ID
// @Description Retrieve a specific approval by its ID
// @Tags Approval
// @Accept json
// @Produce json
// @Param id path int true "Approval ID"
// @Success 200 {object} map[string]interface{} "Approval found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Approval not found"
// @Router /admin/approval/{id} [get]
func (h *ApprovalHandler) GetApproval(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	approval, err := h.approvalService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Approval not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Approval found successfully", approval)
}

// GetApprovalByIdMenu godoc
// @Summary Get approval by ID Menu
// @Description Retrieve a specific approval by its ID Menu
// @Tags Approval
// @Accept json
// @Produce json
// @Param id path int true "menu ID"
// @Success 200 {object} map[string]interface{} "Approval found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID Menu"
// @Failure 404 {object} map[string]interface{} "Not Found: Approval not found"
// @Router /admin/approval/{id}/menu [get]
func (h *ApprovalHandler) GetApprovalByIdMenu(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	approval, err := h.approvalService.GetByIdMenu(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Approval not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Approval found successfully", approval)
}

// CreateApproval godoc
// @Summary Create a new approval
// @Description Create a new approval with the provided details
// @Tags Approval
// @Accept json
// @Produce json
// @Param approval body model.MstApproval true "Approval details"
// @Success 201 {object} map[string]interface{} "Approval created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/approval [post]
func (h *ApprovalHandler) CreateApproval(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusUnauthorized, "Invalid userID"))
	}

	var approval *model.MstApproval
	if err := c.BodyParser(&approval); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, "Failed to parse request body"))
	}

	approval.IDCreatedby = userID
	approval.IDUpdatedby = userID

	if err := h.approvalService.Create(approval); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": approval.ID,
	}

	return pkg.Response(c, fiber.StatusCreated, "Approval created successfully", result)
}

// UpdateApproval godoc
// @Summary Update an existing approval
// @Description Update the details of an existing approval by its ID
// @Tags Approval
// @Accept json
// @Produce json
// @Param id path int true "Approval ID"
// @Param approval body model.MstApproval true "Updated approval details"
// @Success 200 {object} map[string]interface{} "Approval updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 404 {object} map[string]interface{} "Not Found: Approval not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/approval/{id} [put]
func (h *ApprovalHandler) UpdateApproval(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var approval *model.MstApproval
	approval, err = h.approvalService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Approval not found"))
	}

	if err := c.BodyParser(approval); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	approval.ID = uint(ID)
	approval.IDUpdatedby = userID

	err = h.approvalService.Update(approval)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": approval.ID,
	}

	return pkg.Response(c, fiber.StatusOK, "Approval updated successfully", result)
}

// DeleteApproval godoc
// @Summary Delete a approval
// @Description Delete a approval by its ID
// @Tags Approval
// @Param id path int true "Approval ID"
// @Success 200 {object} map[string]interface{} "Approval deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Approval not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/approval/{id} [delete]
func (h *ApprovalHandler) DeleteApproval(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	user, err := h.approvalService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Approval not found"))
	}

	err = h.approvalService.Delete(user)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Approval deleted successfully", nil)
}
