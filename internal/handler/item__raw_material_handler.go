package handler

import (
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type ItemRawMaterialHandler struct {
	itemRawMaterialService *service.ItemRawMaterialService
}

func NewItemRawMaterialHandler(itemRawMaterialService *service.ItemRawMaterialService) *ItemRawMaterialHandler {
	return &ItemRawMaterialHandler{itemRawMaterialService: itemRawMaterialService}
}

// GetItemRawMaterials godoc
// @Summary Get list of Item Raw Materials
// @Tags Item Raw Material
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param rows query int false "Rows per page"
// @Param search query string false "Search"
// @Param categoryCode query string false "Item Category Code"
// @Param sortBy query string false "Sort by field"
// @Param sortDirection query boolean false "true = ASC, false = DESC"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}"
// @Router /general/master/item/generate/raw-material [get]
func (h *ItemRawMaterialHandler) GetItemRawMaterials(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search", "")
	categoryCode := c.Query("categoryCode", "")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection", true)
	offset := (page - 1) * rows

	total, err := h.itemRawMaterialService.GetTotal(search, categoryCode)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	itemRawMaterials, err := h.itemRawMaterialService.GetAll(offset, rows, search, categoryCode, sortBy, sortDirection)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	totalPages := int(math.Ceil(float64(total) / float64(rows)))
	start := offset + 1
	end := int(math.Min(float64(offset+rows), float64(total)))

	result := map[string]interface{}{
		"item_raw_materials": itemRawMaterials,
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

// GetItemRawMaterial godoc
// @Summary Get item raw material by ID
// @Tags Item Raw Material
// @Param id path int true "Item Raw Material ID"
// @Success 200 {object} model.MstItemRawMaterial
// @Failure 404 {object} map[string]interface{}
// @Router /general/master/item/generate/raw-material/{id} [get]
func (h *ItemRawMaterialHandler) GetItemRawMaterial(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, "Invalid ID"))
	}

	itemRawMaterial, err := h.itemRawMaterialService.GetByID(uint(id))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Item Raw Material not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Item Raw Material found successfully", itemRawMaterial)
}

// CreateItemRawMaterial godoc
// @Summary Create new item raw material
// @Tags Item Raw Material
// @Accept json
// @Produce json
// @Param item_raw_material body model.MstItemRawMaterial true "Item Raw Material Body"
// @Success 201 {object} map[string]interface{}
// @Router /general/master/item/generate/raw-material [post]
func (h *ItemRawMaterialHandler) CreateItemRawMaterial(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var itemRawMaterial model.MstItemRawMaterial
	if err := c.BodyParser(&itemRawMaterial); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	itemRawMaterial.IDCreatedby = userID
	itemRawMaterial.IDUpdatedby = userID

	if err := h.itemRawMaterialService.Create(&itemRawMaterial); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusCreated, "Item Raw Material created successfully", map[string]interface{}{"id": itemRawMaterial.ID})
}

// UpdateItemRawMaterial godoc
// @Summary Update item raw material
// @Tags Item Raw Material
// @Param id path int true "Item Raw Material ID"
// @Param item_raw_material body model.MstItemRawMaterial true "Updated Item Raw Material"
// @Success 200 {object} map[string]interface{}
// @Router /general/master/item/generate/raw-material/{id} [put]
func (h *ItemRawMaterialHandler) UpdateItemRawMaterial(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, "Invalid ID"))
	}

	itemRawMaterial, err := h.itemRawMaterialService.GetByID(uint(id))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Item Raw Material not found"))
	}

	if err := c.BodyParser(itemRawMaterial); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	itemRawMaterial.IDUpdatedby = userID

	if err := h.itemRawMaterialService.Update(itemRawMaterial); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Item Raw Material updated successfully", map[string]interface{}{"id": itemRawMaterial.ID})
}

// DeleteItemRawMaterial godoc
// @Summary Delete item raw material
// @Tags Item Raw Material
// @Param id path int true "Item Raw Material ID"
// @Success 200 {object} map[string]interface{}
// @Router /general/master/item/generate/raw-material/{id} [delete]
func (h *ItemRawMaterialHandler) DeleteItemRawMaterial(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, "Invalid ID"))
	}

	itemRawMaterial, err := h.itemRawMaterialService.GetByID(uint(id))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Item Raw Material not found"))
	}

	if err := h.itemRawMaterialService.Delete(itemRawMaterial); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Item Raw Material deleted successfully", nil)
}
