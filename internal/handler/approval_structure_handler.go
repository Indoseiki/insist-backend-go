package handler

import (
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type ApprovalStructureHandler struct {
	approvalStructureService *service.ApprovalStructureService
}

func NewApprovalStructureHandler(approvalStructureService *service.ApprovalStructureService) *ApprovalStructureHandler {
	return &ApprovalStructureHandler{approvalStructureService: approvalStructureService}
}

// GetApprovalStructures godoc
// @Summary Get a list of approval structures
// @Description Retrieves approval structures with pagination and optional search
// @Tags Approval Structure
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering approval structure"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/approval-structure [get]
func (h *ApprovalStructureHandler) GetApprovalStructures(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	total, err := h.approvalStructureService.GetTotal(search)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	approvalStructures, err := h.approvalStructureService.GetAll(offset, rows, search, sortBy, sortDirection)
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
		"items": approvalStructures,
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

	if len(approvalStructures) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetApprovalStructure godoc
// @Summary Get approval structure by ID
// @Description Retrieve a specific approval structure by its ID
// @Tags Approval Structure
// @Accept json
// @Produce json
// @Param id path int true "Approval Structure ID"
// @Success 200 {object} map[string]interface{} "Approval Structure found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Approval Structure not found"
// @Router /admin/approval-structure/{id} [get]
func (h *ApprovalStructureHandler) GetApprovalStructure(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	approvalStructure, err := h.approvalStructureService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Approval Structure not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Approval Structure found successfully", approvalStructure)
}

// GetApprovalStructureByMenu godoc
// @Summary Get approval structure by menu path
// @Description Retrieve approval structure based on the provided menu path
// @Tags Approval Structure
// @Accept json
// @Produce json
// @Param path query string true "Menu Path"
// @Success 200 {object} map[string]interface{} "Approval Structure found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid path"
// @Failure 404 {object} map[string]interface{} "Not Found: Approval Structure not found"
// @Router /admin/approval-structure/menu [get]
func (h *ApprovalStructureHandler) GetApprovalStructureByMenu(c *fiber.Ctx) error {
	path := c.Query("path")
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	approvalStructure, err := h.approvalStructureService.GetAllByMenu(uint(ID), path)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Approval Structure by menu not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Approval Structure by menu found successfully", approvalStructure)
}
