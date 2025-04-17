package handler

import (
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type MaterialHandler struct {
	materialService *service.MaterialService
}

func NewMaterialHandler(materialService *service.MaterialService) *MaterialHandler {
	return &MaterialHandler{materialService: materialService}
}

// GetMaterials godoc
// @Summary Get list of Materials
// @Tags Material
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param rows query int false "Rows per page"
// @Param search query string false "Search"
// @Param sortBy query string false "Sort by field"
// @Param sortDirection query boolean false "true = ASC, false = DESC"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /egd/master/materials [get]
func (h *MaterialHandler) GetMaterials(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search", "")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection", true)
	offset := (page - 1) * rows

	total, err := h.materialService.GetTotal(search)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	materials, err := h.materialService.GetAll(offset, rows, search, sortBy, sortDirection)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	totalPages := int(math.Ceil(float64(total) / float64(rows)))
	start := offset + 1
	end := int(math.Min(float64(offset+rows), float64(total)))

	result := map[string]interface{}{
		"materials": materials,
		"pagination": map[string]interface{}{
			"current_page":  page,
			"total_pages":   totalPages,
			"rows_per_page": rows,
			"total_rows":    total,
			"from":          start,
			"to":            end,
		},
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetMaterial godoc
// @Summary Get Material by ID
// @Tags Material
// @Param id path int true "Material ID"
// @Success 200 {object} model.MstMaterial
// @Failure 404 {object} map[string]interface{}
// @Router /egd/master/material/{id} [get]
func (h *MaterialHandler) GetMaterial(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, "Invalid ID"))
	}

	material, err := h.materialService.GetByID(uint(id))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Material not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Material found successfully", material)
}

// CreateMaterial godoc
// @Summary Create new material
// @Tags Material
// @Accept json
// @Produce json
// @Param material body model.MstMaterial true "Material Body"
// @Success 201 {object} map[string]interface{}
// @Router /egd/master/materials [post]
func (h *MaterialHandler) CreateMaterial(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var material model.MstMaterial
	if err := c.BodyParser(&material); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	material.IDCreatedby = userID
	material.IDUpdatedby = userID

	if err := h.materialService.Create(&material); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusCreated, "Material created successfully", map[string]interface{}{"id": material.ID})
}

// UpdateMaterial godoc
// @Summary Update material
// @Tags Material
// @Param id path int true "Material ID"
// @Param material body model.MstMaterial true "Updated Material"
// @Success 200 {object} map[string]interface{}
// @Router /egd/master/material/{id} [put]
func (h *MaterialHandler) UpdateMaterial(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, "Invalid ID"))
	}

	material, err := h.materialService.GetByID(uint(id))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Material not found"))
	}

	if err := c.BodyParser(material); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	material.IDUpdatedby = userID

	if err := h.materialService.Update(material); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Material updated successfully", map[string]interface{}{"id": material.ID})
}

// DeleteMaterial godoc
// @Summary Delete material
// @Tags Material
// @Param id path int true "Material ID"
// @Success 200 {object} map[string]interface{}
// @Router /egd/master/material/{id} [delete]
func (h *MaterialHandler) DeleteMaterial(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, "Invalid ID"))
	}

	material, err := h.materialService.GetByID(uint(id))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Material not found"))
	}

	if err := h.materialService.Delete(material); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Material deleted successfully", nil)
}
