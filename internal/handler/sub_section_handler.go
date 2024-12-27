package handler

import (
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type SubSectionHandler struct {
	subSectionService *service.SubSectionService
}

func NewSubSectionHandler(subSectionService *service.SubSectionService) *SubSectionHandler {
	return &SubSectionHandler{subSectionService: subSectionService}
}

// GetSubSections godoc
// @Summary Get a list of subSections
// @Description Retrieves subSections with pagination and optional search
// @Tags SubSection
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering subSection"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/sub-section [get]
func (h *SubSectionHandler) GetSubSections(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	total, err := h.subSectionService.GetTotal(search)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	subSections, err := h.subSectionService.GetAll(offset, rows, search, sortBy, sortDirection)
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
		"items": subSections,
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

	if len(subSections) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetSubSection godoc
// @Summary Get subSection by ID
// @Description Retrieve a specific subSection by its ID
// @Tags SubSection
// @Accept json
// @Produce json
// @Param id path int true "SubSection ID"
// @Success 200 {object} map[string]interface{} "SubSection found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: SubSection not found"
// @Router /admin/master/sub-section/{id} [get]
func (h *SubSectionHandler) GetSubSection(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	subSection, err := h.subSectionService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "SubSection not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "SubSection found successfully", subSection)
}

// CreateSubSection godoc
// @Summary Create a new subSection
// @Description Create a new subSection with the provided details
// @Tags SubSection
// @Accept json
// @Produce json
// @Param subSection body model.MstSubSection true "SubSection details"
// @Success 201 {object} map[string]interface{} "SubSection created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/sub-section [post]
func (h *SubSectionHandler) CreateSubSection(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var subSection model.MstSubSection
	if err := c.BodyParser(&subSection); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	subSection.IDCreatedby = userID
	subSection.IDUpdatedby = userID

	err := h.subSectionService.Create(&subSection)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": subSection.ID,
	}

	return pkg.Response(c, fiber.StatusCreated, "SubSection created successfully", result)
}

// UpdateSubSection godoc
// @Summary Update an existing subSection
// @Description Update the details of an existing subSection by its ID
// @Tags SubSection
// @Accept json
// @Produce json
// @Param id path int true "SubSection ID"
// @Param subSection body model.MstSubSection true "Updated subSection details"
// @Success 200 {object} map[string]interface{} "SubSection updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 404 {object} map[string]interface{} "Not Found: SubSection not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/sub-section/{id} [put]
func (h *SubSectionHandler) UpdateSubSection(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var subSection *model.MstSubSection
	subSection, err = h.subSectionService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "SubSection not found"))
	}

	if err := c.BodyParser(subSection); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	subSection.ID = uint(ID)
	subSection.IDUpdatedby = userID

	err = h.subSectionService.Update(subSection)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": subSection.ID,
	}

	return pkg.Response(c, fiber.StatusOK, "SubSection updated successfully", result)
}

// DeleteSubSection godoc
// @Summary Delete a subSection
// @Description Delete a subSection by its ID
// @Tags SubSection
// @Param id path int true "SubSection ID"
// @Success 200 {object} map[string]interface{} "SubSection deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: SubSection not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/sub-section/{id} [delete]
func (h *SubSectionHandler) DeleteSubSection(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	user, err := h.subSectionService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "SubSection not found"))
	}

	err = h.subSectionService.Delete(user)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "SubSection deleted successfully", nil)
}
