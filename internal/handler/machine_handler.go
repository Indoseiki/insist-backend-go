package handler

import (
	"insist-backend-golang/internal/dto"
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type MachineHandler struct {
	machineService *service.MachineService
}

func NewMachineHandler(machineService *service.MachineService) *MachineHandler {
	return &MachineHandler{machineService: machineService}
}

// GetMachines godoc
// @Summary Get a list of machines
// @Description Retrieves machines with pagination and optional search
// @Tags Machine
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering machine"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /mnt/master/machine [get]
func (h *MachineHandler) GetMachines(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	total, err := h.machineService.GetTotal(search)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	machines, err := h.machineService.GetAll(offset, rows, search, sortBy, sortDirection)
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
		"items": machines,
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

	if len(machines) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetMachine godoc
// @Summary Get machine by ID
// @Description Retrieve a specific machine by its ID
// @Tags Machine
// @Accept json
// @Produce json
// @Param id path int true "Machine ID"
// @Success 200 {object} map[string]interface{} "Machine found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Machine not found"
// @Router /mnt/master/machine/{id} [get]
func (h *MachineHandler) GetMachine(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	machine, err := h.machineService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Machine not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Machine found successfully", machine)
}

// CreateMachine godoc
// @Summary Create a new machine
// @Description Create a new machine with the provided details
// @Tags Machine
// @Accept json
// @Produce json
// @Param machine body dto.CreateMachinePayload true "Machine and its details"
// @Success 201 {object} map[string]interface{} "Machine created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /mnt/master/machine [post]
func (h *MachineHandler) CreateMachine(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var payload dto.CreateMachinePayload

	if err := c.BodyParser(&payload); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	payload.Machine.IDCreatedby = userID
	payload.Machine.IDUpdatedby = userID

	tx := h.machineService.BeginTx()

	if err := tx.Create(&payload.Machine).Error; err != nil {
		tx.Rollback()
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	payload.MachineDetail.RevNo = 0
	payload.MachineDetail.IDMachine = payload.Machine.ID
	payload.MachineStatus.IDMachine = payload.Machine.ID

	payload.MachineDetail.IDCreatedby = userID
	payload.MachineDetail.IDUpdatedby = userID
	payload.MachineStatus.IDCreatedby = userID
	payload.MachineStatus.IDUpdatedby = userID

	if err := tx.Create(&payload.MachineDetail).Error; err != nil {
		tx.Rollback()
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	if err := tx.Create(&payload.MachineStatus).Error; err != nil {
		tx.Rollback()
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	tx.Commit()

	result := map[string]interface{}{
		"id":             payload.Machine.ID,
		"machine_detail": payload.MachineDetail,
		"machine_status": payload.MachineStatus,
	}

	return pkg.Response(c, fiber.StatusCreated, "Machine created successfully", result)
}

// UpdateMachine godoc
// @Summary Update an existing machine
// @Description Update the details of an existing machine by its ID
// @Tags Machine
// @Accept json
// @Produce json
// @Param id path int true "Machine ID"
// @Param machine body model.MstMachine true "Updated machine details"
// @Success 200 {object} map[string]interface{} "Machine updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 404 {object} map[string]interface{} "Not Found: Machine not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /mnt/master/machine/{id} [put]
func (h *MachineHandler) UpdateMachine(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	detailID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	machineDetail, err := h.machineService.GetDetailByID(uint(detailID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Detail machine not found"))
	}

	if err := c.BodyParser(machineDetail); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	machineDetail.IDUpdatedby = userID

	err = h.machineService.UpdateDetail(machineDetail)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": machineDetail.ID,
	}

	return pkg.Response(c, fiber.StatusOK, "Machine detail updated successfully", result)
}

// DeleteMachine godoc
// @Summary Delete a machine
// @Description Delete a machine by its ID
// @Tags Machine
// @Param id path int true "Machine ID"
// @Success 200 {object} map[string]interface{} "Machine deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Machine not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /mnt/master/machine/{id} [delete]
func (h *MachineHandler) DeleteMachine(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	machine, err := h.machineService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, "Machine not found"))
	}

	machineDetails, err := h.machineService.GetDetailsByMachineID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, "Detail machine not found"))
	}

	if len(machineDetails) > 1 {
		err = h.machineService.DeleteDetail(uint(machine.DetailID))
		if err != nil {
			return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
		}
	} else if len(machineDetails) == 1 {
		err = h.machineService.DeleteDetail(uint(machine.DetailID))
		if err != nil {
			return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
		}

		err = h.machineService.DeleteStatus(uint(ID))
		if err != nil {
			return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
		}

		err = h.machineService.Delete(uint(ID))
		if err != nil {
			return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
		}
	} else {
		err = h.machineService.DeleteStatus(uint(ID))
		if err != nil {
			return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
		}

		err = h.machineService.Delete(uint(ID))
		if err != nil {
			return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
		}
	}

	return pkg.Response(c, fiber.StatusOK, "Machine deleted successfully", nil)
}

// RevisionMachine godoc
// @Summary Create a new revision of an existing machine
// @Description This endpoint creates a new revision of an existing machine by copying its details and updating the revision number.
// @Tags Machine
// @Accept json
// @Produce json
// @Param id path int true "Detail Machine ID"
// @Success 201 {object} map[string]interface{} "Machine revision created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid machine ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Machine not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /mnt/master/machine/{id}/revision [put]
func (h *MachineHandler) RevisionMachine(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	existingMachineDetail, err := h.machineService.GetDetailByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Machine not found"))
	}

	newMachineDetail := *existingMachineDetail
	newMachineDetail.ID = 0
	newMachineDetail.RevNo += 1
	newMachineDetail.IDCreatedby = userID
	newMachineDetail.IDUpdatedby = userID

	tx := h.machineService.BeginTx()

	if err := tx.Create(&newMachineDetail).Error; err != nil {
		tx.Rollback()
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	tx.Commit()

	result := map[string]interface{}{
		"id":     newMachineDetail.ID,
		"rev_no": newMachineDetail.RevNo,
	}

	return pkg.Response(c, fiber.StatusCreated, "Machine revision created successfully", result)
}

// GetMachineDetails godoc
// @Summary Retrieve machine details with pagination, sorting, and filtering
// @Description Fetches a list of machine details with optional sorting, and pagination.
// @Tags Machine
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param rows query int false "Number of rows per page (default: 20)"
// @Param sortBy query string false "Column name to sort by"
// @Param sortDirection query boolean false "Sorting direction: true for ascending, false for descending"
// @Success 200 {object} map[string]interface{} "Machine details retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid query parameters"
// @Failure 404 {object} map[string]interface{} "Not Found: No machine details found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /mnt/master/machine/{id}/detail [get]
func (h *MachineHandler) GetMachineDetails(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	total, err := h.machineService.GetTotalDetail(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	machineDetails, err := h.machineService.GetAllDetail(offset, rows, sortBy, sortDirection, uint(ID))
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
		"items": machineDetails,
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

	if len(machineDetails) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetMachineStatus godoc
// @Summary Retrieve machine status details with pagination, sorting, and filtering
// @Description Fetches a list of machine status details for a specific machine ID with optional sorting and pagination.
// @Tags Machine
// @Accept json
// @Produce json
// @Param id path int true "Machine ID"
// @Param page query int false "Page number (default: 1)"
// @Param rows query int false "Number of rows per page (default: 20)"
// @Param sortBy query string false "Column name to sort by (default: detail_updated_at)"
// @Param sortDirection query boolean false "Sorting direction: false for ascending, true for descending (default: false)"
// @Success 200 {object} map[string]interface{} "Machine status details retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid query parameters"
// @Failure 404 {object} map[string]interface{} "Not Found: No machine status details found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /mnt/master/machine/{id}/status [get]
func (h *MachineHandler) GetMachineStatus(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	total, err := h.machineService.GetTotalStatus(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	machineStatus, err := h.machineService.GetAllStatus(offset, rows, sortBy, sortDirection, uint(ID))
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
		"items": machineStatus,
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

	if len(machineStatus) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// CreateStatusMachine godoc
// @Summary Create a new machine status
// @Description Adds a new status entry for a specific machine ID.
// @Tags Machine
// @Accept json
// @Produce json
// @Param id path int true "Machine ID"
// @Param body body model.MstMachineStatus true "Machine status details"
// @Success 201 {object} map[string]interface{} "Machine status created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /mnt/master/machine/{id}/status [put]
func (h *MachineHandler) CreateStatusMachine(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var machineStatus model.MstMachineStatus
	if err := c.BodyParser(&machineStatus); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	machineStatus.IDMachine = uint(ID)
	machineStatus.IDCreatedby = userID
	machineStatus.IDUpdatedby = userID

	err = h.machineService.CreateStatus(&machineStatus)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": machineStatus.ID,
	}

	return pkg.Response(c, fiber.StatusCreated, "Machine status created successfully", result)
}
