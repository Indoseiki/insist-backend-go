package handler

import (
	"fmt"
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"
	"strings"

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
	action := c.Query("action", "")
	isSuccess := c.Query("isSuccess", "")
	rangeDate := c.Query("rangeDate", "")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	arrayDate := strings.Split(rangeDate, "~")
	fmt.Println(arrayDate)

	total, err := h.ActivityLogService.GetTotal(search, action, isSuccess, arrayDate)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	ActivityLogs, err := h.ActivityLogService.GetAll(offset, rows, search, action, isSuccess, sortBy, sortDirection, arrayDate)
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
		"items": ActivityLogs,
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
// @Param ActivityLog body model.ActivityLog true "Activity Log details"
// @Success 201 {object} map[string]interface{} "Activity Log created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/ActivityLog [post]
func (h *ActivityLogHandler) CreateActivityLog(c *fiber.Ctx) error {
	var ActivityLog model.ActivityLog
	if err := c.BodyParser(&ActivityLog); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	user, err := h.ActivityLogService.GetByUsername(ActivityLog.Username)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusOK, err.Error()))
	}

	ActivityLog.IDUser = user.ID
	ActivityLog.IPAddress = c.IP()

	err = h.ActivityLogService.Create(&ActivityLog)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": ActivityLog.ID,
	}

	return pkg.Response(c, fiber.StatusCreated, "Activity Log created successfully", result)
}
