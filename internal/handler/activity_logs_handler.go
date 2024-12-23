package handler

import (
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type ActivityLogHandler struct {
	ActivityLogService *service.ActivityLogService
}

func NewActivityLogHandler(ActivityLogService *service.ActivityLogService) *ActivityLogHandler {
	return &ActivityLogHandler{ActivityLogService: ActivityLogService}
}

// GetActivityLogs godoc
// @Summary Get a list of Activity Logs
// @Description Retrieves Activity Logs with pagination and optional search
// @Tags Activity Log
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering Activity Log"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/ActivityLog [get]
func (h *ActivityLogHandler) GetActivityLogs(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	total, err := h.ActivityLogService.GetTotal(search)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	ActivityLogs, err := h.ActivityLogService.GetAll(offset, rows, search, sortBy, sortDirection)
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
		"items": ActivityLogs,
		"pagination": map[string]interface{}{
			"current_page":  page,
			"next_page":     nextPage,
			"total_pages":   totalPages,
			"rows_per_page": rows,
			"total_rows":    total,
		},
	}

	if len(ActivityLogs) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetActivityLog godoc
// @Summary Get Activity Log by ID
// @Description Retrieve a specific Activity Log by its ID
// @Tags Activity Log
// @Accept json
// @Produce json
// @Param id path int true "Activity Log ID"
// @Success 200 {object} map[string]interface{} "Activity Log found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Activity Log not found"
// @Router /admin/master/ActivityLog/{id} [get]
func (h *ActivityLogHandler) GetActivityLog(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	ActivityLog, err := h.ActivityLogService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Activity Log not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Activity Log found successfully", ActivityLog)
}

// CreateActivityLog godoc
// @Summary Create a new Activity Log
// @Description Create a new Activity Log with the provided details
// @Tags Activity Log
// @Accept json
// @Produce json
// @Param Activity Log body model.ActivityLog true "Activity Log details"
// @Success 201 {object} map[string]interface{} "Activity Log created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/ActivityLog [post]
func (h *ActivityLogHandler) CreateActivityLog(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var ActivityLog model.ActivityLog
	if err := c.BodyParser(&ActivityLog); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	ActivityLog.IDUser = userID
	ActivityLog.IPAddress = c.IP()

	err := h.ActivityLogService.Create(&ActivityLog)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": ActivityLog.ID,
	}

	return pkg.Response(c, fiber.StatusCreated, "Activity Log created successfully", result)
}

// UpdateActivityLog godoc
// @Summary Update an existing Activity Log
// @Description Update the details of an existing Activity Log by its ID
// @Tags Activity Log
// @Accept json
// @Produce json
// @Param id path int true "Activity Log ID"
// @Param Activity Log body model.ActivityLog true "Updated Activity Log details"
// @Success 200 {object} map[string]interface{} "Activity Log updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 404 {object} map[string]interface{} "Not Found: Activity Log not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/ActivityLog/{id} [put]
func (h *ActivityLogHandler) UpdateActivityLog(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var ActivityLog *model.ActivityLog
	ActivityLog, err = h.ActivityLogService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Activity Log not found"))
	}

	if err := c.BodyParser(ActivityLog); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	ActivityLog.ID = uint(ID)

	err = h.ActivityLogService.Update(ActivityLog)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": ActivityLog.ID,
	}

	return pkg.Response(c, fiber.StatusOK, "Activity Log updated successfully", result)
}

// DeleteActivityLog godoc
// @Summary Delete a Activity Log
// @Description Delete a Activity Log by its ID
// @Tags Activity Log
// @Param id path int true "Activity Log ID"
// @Success 200 {object} map[string]interface{} "Activity Log deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Activity Log not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/ActivityLog/{id} [delete]
func (h *ActivityLogHandler) DeleteActivityLog(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	user, err := h.ActivityLogService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Activity Log not found"))
	}

	err = h.ActivityLogService.Delete(user)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Activity Log deleted successfully", nil)
}
