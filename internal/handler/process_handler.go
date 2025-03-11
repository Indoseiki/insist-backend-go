package handler

import (
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type ProcessHandler struct {
	processService *service.ProcessService
}

func NewProcessHandler(processService *service.ProcessService) *ProcessHandler {
	return &ProcessHandler{processService: processService}
}

// GetProcesss godoc
// @Summary Get a list of processs
// @Description Retrieves processs with pagination and optional search
// @Tags Process
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering process"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /egd/master/process [get]
func (h *ProcessHandler) GetProcesses(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	total, err := h.processService.GetTotal(search)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	process, err := h.processService.GetAll(offset, rows, search, sortBy, sortDirection)
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
		"items": process,
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

	if len(process) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetProcess godoc
// @Summary Get process by ID
// @Description Retrieve a specific process by its ID
// @Tags Process
// @Accept json
// @Produce json
// @Param id path int true "Process ID"
// @Success 200 {object} map[string]interface{} "Process found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Process not found"
// @Router /egd/master/process/{id} [get]
func (h *ProcessHandler) GetProcess(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	process, err := h.processService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Process not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Process found successfully", process)
}

// CreateProcess godoc
// @Summary Create a new process
// @Description Create a new process with the provided details
// @Tags Process
// @Accept json
// @Produce json
// @Param process body model.MstProcess true "Process details"
// @Success 201 {object} map[string]interface{} "Process created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /egd/master/process [post]
func (h *ProcessHandler) CreateProcess(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var process model.MstProcess
	if err := c.BodyParser(&process); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	process.IDCreatedby = userID
	process.IDUpdatedby = userID

	err := h.processService.Create(&process)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": process.ID,
	}

	return pkg.Response(c, fiber.StatusCreated, "Process created successfully", result)
}

// UpdateProcess godoc
// @Summary Update an existing process
// @Description Update the details of an existing process by its ID
// @Tags Process
// @Accept json
// @Produce json
// @Param id path int true "Process ID"
// @Param process body model.MstProcess true "Updated process details"
// @Success 200 {object} map[string]interface{} "Process updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 404 {object} map[string]interface{} "Not Found: Process not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /egd/master/process/{id} [put]
func (h *ProcessHandler) UpdateProcess(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var process *model.MstProcess
	process, err = h.processService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Process not found"))
	}

	if err := c.BodyParser(process); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	process.ID = uint(ID)
	process.IDUpdatedby = userID

	err = h.processService.Update(process)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": process.ID,
	}

	return pkg.Response(c, fiber.StatusOK, "Process updated successfully", result)
}

// DeleteProcess godoc
// @Summary Delete a process
// @Description Delete a process by its ID
// @Tags Process
// @Param id path int true "Process ID"
// @Success 200 {object} map[string]interface{} "Process deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Process not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /egd/master/process/{id} [delete]
func (h *ProcessHandler) DeleteProcess(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	user, err := h.processService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Process not found"))
	}

	err = h.processService.Delete(user)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Process deleted successfully", nil)
}
