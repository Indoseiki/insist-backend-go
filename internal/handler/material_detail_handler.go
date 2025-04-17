package handler

import (
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type MaterialDetailHandler struct {
	materialDetailService *service.MaterialDetailService
}

func NewMaterialDetailHandler(materialDetailService *service.MaterialDetailService) *MaterialDetailHandler {
	return &MaterialDetailHandler{materialDetailService: materialDetailService}
}

// GetMaterialDetails godoc
// @Summary Get list of Material Details
// @Tags Material Detail
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param rows query int false "Rows per page"
// @Param search query string false "Search"
// @Param idMaterial query int false "Material ID"
// @Param sortBy query string false "Sort by field"
// @Param sortDirection query boolean false "true = ASC, false = DESC"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /egd/master/material-detail [get]
func (h *MaterialDetailHandler) GetMaterialDetails(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search", "")
	idMaterial := c.QueryInt("idMaterial", 0)
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection", true)
	offset := (page - 1) * rows

	total, err := h.materialDetailService.GetTotal(search, uint(idMaterial))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	materialDetails, err := h.materialDetailService.GetAll(offset, rows, search, sortBy, sortDirection, uint(idMaterial))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	totalPages := int(math.Ceil(float64(total) / float64(rows)))
	start := offset + 1
	end := int(math.Min(float64(offset+rows), float64(total)))

	result := map[string]interface{}{
		"materialDetails": materialDetails,
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

// GetMaterialDetail godoc
// @Summary Get Material Detail by ID
// @Tags Material Detail
// @Param id path int true "Material Detail ID"
// @Success 200 {object} model.MstMaterialDetail
// @Failure 404 {object} map[string]interface{}
// @Router /egd/master/material-detail/{id} [get]
func (h *MaterialDetailHandler) GetMaterialDetail(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, "Invalid ID"))
	}

	materialDetail, err := h.materialDetailService.GetByID(uint(id))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Material Detail not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Material Detail found successfully", materialDetail)
}

// CreateMaterialDetail godoc
// @Summary Create new material detail
// @Tags Material Detail
// @Accept json
// @Produce json
// @Param materialDetail body model.MstMaterialDetail true "Material Detail Body"
// @Success 201 {object} map[string]interface{}
// @Router /egd/master/material-detail [post]
func (h *MaterialDetailHandler) CreateMaterialDetail(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var materialDetail model.MstMaterialDetail
	if err := c.BodyParser(&materialDetail); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	materialDetail.IDCreatedby = userID
	materialDetail.IDUpdatedby = userID

	if err := h.materialDetailService.Create(&materialDetail); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusCreated, "Material Detail created successfully", map[string]interface{}{"id": materialDetail.ID})
}

// UpdateMaterialDetail godoc
// @Summary Update material detail
// @Tags Material Detail
// @Param id path int true "Material Detail ID"
// @Param materialDetail body model.MstMaterialDetail true "Updated Material Detail"
// @Success 200 {object} map[string]interface{}
// @Router /egd/master/material-detail/{id} [put]
func (h *MaterialDetailHandler) UpdateMaterialDetail(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, "Invalid ID"))
	}

	materialDetail, err := h.materialDetailService.GetByID(uint(id))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Material Detail not found"))
	}

	if err := c.BodyParser(materialDetail); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	materialDetail.IDUpdatedby = userID

	if err := h.materialDetailService.Update(materialDetail); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Material Detail updated successfully", map[string]interface{}{"id": materialDetail.ID})
}

// DeleteMaterialDetail godoc
// @Summary Delete material detail
// @Tags Material Detail
// @Param id path int true "Material Detail ID"
// @Success 200 {object} map[string]interface{}
// @Router /egd/master/material-detail/{id} [delete]
func (h *MaterialDetailHandler) DeleteMaterialDetail(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, "Invalid ID"))
	}

	materialDetail, err := h.materialDetailService.GetByID(uint(id))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Material Detail not found"))
	}

	if err := h.materialDetailService.Delete(materialDetail); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Material Detail deleted successfully", nil)
}
