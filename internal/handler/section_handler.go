package handler

import (
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type SectionHandler struct {
	sectionService *service.SectionService
}

func NewSectionHandler(sectionService *service.SectionService) *SectionHandler {
	return &SectionHandler{sectionService: sectionService}
}

// GetSections godoc
// @Summary Get a list of sections
// @Description Retrieves sections with pagination and optional search
// @Tags Section
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering section"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/section [get]
func (h *SectionHandler) GetSections(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	idFCS := c.QueryInt("id_fcs")
	search := c.Query("search")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	total, err := h.sectionService.GetTotal(search, uint(idFCS))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	sections, err := h.sectionService.GetAll(offset, rows, search, sortBy, sortDirection, uint(idFCS))
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
		"items": sections,
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

	if len(sections) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetSection godoc
// @Summary Get section by ID
// @Description Retrieve a specific section by its ID
// @Tags Section
// @Accept json
// @Produce json
// @Param id path int true "Section ID"
// @Success 200 {object} map[string]interface{} "Section found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Section not found"
// @Router /admin/master/section/{id} [get]
func (h *SectionHandler) GetSection(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	section, err := h.sectionService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Section not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Section found successfully", section)
}

// CreateSection godoc
// @Summary Create a new section
// @Description Create a new section with the provided details
// @Tags Section
// @Accept json
// @Produce json
// @Param section body model.MstSection true "Section details"
// @Success 201 {object} map[string]interface{} "Section created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/section [post]
func (h *SectionHandler) CreateSection(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var section model.MstSection
	if err := c.BodyParser(&section); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	section.IDCreatedby = userID
	section.IDUpdatedby = userID

	err := h.sectionService.Create(&section)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": section.ID,
	}

	return pkg.Response(c, fiber.StatusCreated, "Section created successfully", result)
}

// UpdateSection godoc
// @Summary Update an existing section
// @Description Update the details of an existing section by its ID
// @Tags Section
// @Accept json
// @Produce json
// @Param id path int true "Section ID"
// @Param section body model.MstSection true "Updated section details"
// @Success 200 {object} map[string]interface{} "Section updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 404 {object} map[string]interface{} "Not Found: Section not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/section/{id} [put]
func (h *SectionHandler) UpdateSection(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var section *model.MstSection
	section, err = h.sectionService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Section not found"))
	}

	if err := c.BodyParser(section); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	section.ID = uint(ID)
	section.IDUpdatedby = userID

	err = h.sectionService.Update(section)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": section.ID,
	}

	return pkg.Response(c, fiber.StatusOK, "Section updated successfully", result)
}

// DeleteSection godoc
// @Summary Delete a section
// @Description Delete a section by its ID
// @Tags Section
// @Param id path int true "Section ID"
// @Success 200 {object} map[string]interface{} "Section deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Section not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/section/{id} [delete]
func (h *SectionHandler) DeleteSection(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	user, err := h.sectionService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Section not found"))
	}

	err = h.sectionService.Delete(user)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Section deleted successfully", nil)
}
