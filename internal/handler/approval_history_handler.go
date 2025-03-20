package handler

import (
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type ApprovalHistoryHandler struct {
	approvalHistoryService *service.ApprovalHistoryService
}

func NewApprovalHistoryHandler(approvalHistoryService *service.ApprovalHistoryService) *ApprovalHistoryHandler {
	return &ApprovalHistoryHandler{approvalHistoryService: approvalHistoryService}
}

// GetApprovalHistories godoc
// @Summary Get a list of approvalHistorys
// @Description Retrieves approvalHistorys with pagination and optional search
// @Tags ApprovalHistory
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering approvalHistory"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/approval-history [get]
func (h *ApprovalHistoryHandler) GetApprovalHistories(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	total, err := h.approvalHistoryService.GetTotal(search)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	approvalHistorys, err := h.approvalHistoryService.GetAll(offset, rows, search, sortBy, sortDirection)
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
		"items": approvalHistorys,
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

	if len(approvalHistorys) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetApprovalHistory godoc
// @Summary Get approvalHistory by ID
// @Description Retrieve a specific approvalHistory by its ID
// @Tags ApprovalHistory
// @Accept json
// @Produce json
// @Param id path int true "Approval History ID"
// @Success 200 {object} map[string]interface{} "Approval History found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: ApprovalHistory not found"
// @Router /admin/approval-history/{id} [get]
func (h *ApprovalHistoryHandler) GetApprovalHistory(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	approvalHistory, err := h.approvalHistoryService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Approval History not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Approval History found successfully", approvalHistory)
}

// CreateApprovalHistory godoc
// @Summary Create a new approvalHistory
// @Description Create a new approvalHistory with the provided details
// @Tags ApprovalHistory
// @Accept json
// @Produce json
// @Param approvalHistory body model.ApprovalHistory true "Approval History details"
// @Success 201 {object} map[string]interface{} "Approval History created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/approval-history [post]
func (h *ApprovalHistoryHandler) CreateApprovalHistory(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var approvalHistory model.ApprovalHistory
	if err := c.BodyParser(&approvalHistory); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	approvalHistory.IDCreatedby = userID

	err := h.approvalHistoryService.Create(&approvalHistory)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": approvalHistory.ID,
	}

	return pkg.Response(c, fiber.StatusCreated, "Approval History created successfully", result)
}

// UpdateApprovalHistory godoc
// @Summary Update an existing approvalHistory
// @Description Update the details of an existing approvalHistory by its ID
// @Tags ApprovalHistory
// @Accept json
// @Produce json
// @Param id path int true "Approval History ID"
// @Param approvalHistory body model.ApprovalHistory true "Updated approvalHistory details"
// @Success 200 {object} map[string]interface{} "Approval History updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 404 {object} map[string]interface{} "Not Found: ApprovalHistory not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/approval-history/{id} [put]
func (h *ApprovalHistoryHandler) UpdateApprovalHistory(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var approvalHistory *model.ApprovalHistory
	approvalHistory, err = h.approvalHistoryService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Approval History not found"))
	}

	if err := c.BodyParser(approvalHistory); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	approvalHistory.ID = uint(ID)
	approvalHistory.IDCreatedby = userID

	err = h.approvalHistoryService.Update(approvalHistory)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": approvalHistory.ID,
	}

	return pkg.Response(c, fiber.StatusOK, "Approval History updated successfully", result)
}

// DeleteApprovalHistory godoc
// @Summary Delete a approvalHistory
// @Description Delete a approvalHistory by its ID
// @Tags ApprovalHistory
// @Param id path int true "Approval History ID"
// @Success 200 {object} map[string]interface{} "Approval History deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: ApprovalHistory not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/approval-history/{id} [delete]
func (h *ApprovalHistoryHandler) DeleteApprovalHistory(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	user, err := h.approvalHistoryService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Approval History not found"))
	}

	err = h.approvalHistoryService.Delete(user)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Approval History deleted successfully", nil)
}

// GetApprovalNotifications godoc
// @Summary Get approval notifications for the logged-in user
// @Description Retrieve approval notifications based on the authenticated user's ID
// @Tags Approval Notifications
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Approval Notifications found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: Approval Notifications not found"
// @Router /admin/approval-history/notifications [get]
func (h *ApprovalHistoryHandler) GetApprovalNotifications(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	approvalNotifications, err := h.approvalHistoryService.GetNotification(userID)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Approval Notification not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Approval Notification found successfully", approvalNotifications)
}

// GetAllByRefID godoc
// @Summary Get approval histories by reference ID
// @Description Retrieve approval histories based on the provided reference ID (as a path parameter) and reference table
// @Tags Approval Histories
// @Accept json
// @Produce json
// @Param id path int true "Reference ID"
// @Param ref_table query string false "Reference Table"
// @Success 200 {object} map[string]interface{} "Approval Histories found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID parameter"
// @Failure 404 {object} map[string]interface{} "Not Found: Approval Histories not found"
// @Router /admin/approval-history/{id}/ref [get]
func (h *ApprovalHistoryHandler) GetAllByRefID(c *fiber.Ctx) error {
	refTable := c.Query("ref_table")
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	approvalHistories, err := h.approvalHistoryService.GetAllByRefID(uint(ID), refTable)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Approval Histories not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Approval Histories found successfully", approvalHistories)
}
