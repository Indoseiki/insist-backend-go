package handler

import (
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type UoMHandler struct {
	uomService *service.UoMService
}

func NewUoMHandler(uomService *service.UoMService) *UoMHandler {
	return &UoMHandler{uomService: uomService}
}

// GetUoMs godoc
// @Summary Get a list of uoms
// @Description Retrieves uoms with pagination and optional search
// @Tags UoM
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering uom"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /egd/master/uom [get]
func (h *UoMHandler) GetUoMs(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	total, err := h.uomService.GetTotal(search)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	uoms, err := h.uomService.GetAll(offset, rows, search, sortBy, sortDirection)
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
		"items": uoms,
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

	if len(uoms) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetUoM godoc
// @Summary Get uom by ID
// @Description Retrieve a specific uom by its ID
// @Tags UoM
// @Accept json
// @Produce json
// @Param id path int true "UoM ID"
// @Success 200 {object} map[string]interface{} "UoM found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: UoM not found"
// @Router /egd/master/uom/{id} [get]
func (h *UoMHandler) GetUoM(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	uom, err := h.uomService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "UoM not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "UoM found successfully", uom)
}

// CreateUoM godoc
// @Summary Create a new uom
// @Description Create a new uom with the provided details
// @Tags UoM
// @Accept json
// @Produce json
// @Param uom body model.MstUoms true "UoM details"
// @Success 201 {object} map[string]interface{} "UoM created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /egd/master/uom [post]
func (h *UoMHandler) CreateUoM(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var uom model.MstUoms
	if err := c.BodyParser(&uom); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	uom.IDCreatedby = userID
	uom.IDUpdatedby = userID

	err := h.uomService.Create(&uom)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": uom.ID,
	}

	return pkg.Response(c, fiber.StatusCreated, "UoM created successfully", result)
}

// UpdateUoM godoc
// @Summary Update an existing uom
// @Description Update the details of an existing uom by its ID
// @Tags UoM
// @Accept json
// @Produce json
// @Param id path int true "UoM ID"
// @Param uom body model.MstUoms true "Updated uom details"
// @Success 200 {object} map[string]interface{} "UoM updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 404 {object} map[string]interface{} "Not Found: UoM not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /egd/master/uom/{id} [put]
func (h *UoMHandler) UpdateUoM(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var uom *model.MstUoms
	uom, err = h.uomService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "UoM not found"))
	}

	if err := c.BodyParser(uom); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	uom.ID = uint(ID)
	uom.IDUpdatedby = userID

	err = h.uomService.Update(uom)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": uom.ID,
	}

	return pkg.Response(c, fiber.StatusOK, "UoM updated successfully", result)
}

// DeleteUoM godoc
// @Summary Delete a uom
// @Description Delete a uom by its ID
// @Tags UoM
// @Param id path int true "UoM ID"
// @Success 200 {object} map[string]interface{} "UoM deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: UoM not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /egd/master/uom/{id} [delete]
func (h *UoMHandler) DeleteUoM(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	user, err := h.uomService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "UoM not found"))
	}

	err = h.uomService.Delete(user)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "UoM deleted successfully", nil)
}
