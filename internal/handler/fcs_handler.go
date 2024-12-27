package handler

import (
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type FCSHandler struct {
	fcsService *service.FCSService
}

func NewFCSHandler(fcsService *service.FCSService) *FCSHandler {
	return &FCSHandler{fcsService: fcsService}
}

// GetFCSs godoc
// @Summary Get a list of fcs
// @Description Retrieves fcs with pagination and optional search
// @Tags FCS
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering fcs"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/fcs [get]
func (h *FCSHandler) GetFCSs(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	total, err := h.fcsService.GetTotal(search)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	fcs, err := h.fcsService.GetAll(offset, rows, search, sortBy, sortDirection)
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
		"items": fcs,
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

	if len(fcs) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetFCS godoc
// @Summary Get fcs by ID
// @Description Retrieve a specific fcs by its ID
// @Tags FCS
// @Accept json
// @Produce json
// @Param id path int true "FCS ID"
// @Success 200 {object} map[string]interface{} "FCS found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: FCS not found"
// @Router /admin/master/fcs/{id} [get]
func (h *FCSHandler) GetFCS(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	fcs, err := h.fcsService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "FCS not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "FCS found successfully", fcs)
}

// CreateFCS godoc
// @Summary Create a new fcs
// @Description Create a new fcs with the provided details
// @Tags FCS
// @Accept json
// @Produce json
// @Param fcs body model.MstFCS true "FCS details"
// @Success 201 {object} map[string]interface{} "FCS created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/fcs [post]
func (h *FCSHandler) CreateFCS(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var fcs model.MstFCS
	if err := c.BodyParser(&fcs); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	fcs.IDCreatedby = userID
	fcs.IDUpdatedby = userID

	err := h.fcsService.Create(&fcs)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": fcs.ID,
	}

	return pkg.Response(c, fiber.StatusCreated, "FCS created successfully", result)
}

// UpdateFCS godoc
// @Summary Update an existing fcs
// @Description Update the details of an existing fcs by its ID
// @Tags FCS
// @Accept json
// @Produce json
// @Param id path int true "FCS ID"
// @Param fcs body model.MstFCS true "Updated fcs details"
// @Success 200 {object} map[string]interface{} "FCS updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 404 {object} map[string]interface{} "Not Found: FCS not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/fcs/{id} [put]
func (h *FCSHandler) UpdateFCS(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var fcs *model.MstFCS
	fcs, err = h.fcsService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "FCS not found"))
	}

	if err := c.BodyParser(fcs); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	fcs.ID = uint(ID)
	fcs.IDUpdatedby = userID

	err = h.fcsService.Update(fcs)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": fcs.ID,
	}

	return pkg.Response(c, fiber.StatusOK, "FCS updated successfully", result)
}

// DeleteFCS godoc
// @Summary Delete a fcs
// @Description Delete a fcs by its ID
// @Tags FCS
// @Param id path int true "FCS ID"
// @Success 200 {object} map[string]interface{} "FCS deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: FCS not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/fcs/{id} [delete]
func (h *FCSHandler) DeleteFCS(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	user, err := h.fcsService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "FCS not found"))
	}

	err = h.fcsService.Delete(user)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "FCS deleted successfully", nil)
}
