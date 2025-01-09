package handler

import (
	"insist-backend-golang/internal/dto"
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type FCSBuildingHandler struct {
	userRoleService *service.FCSBuildingService
}

func NewFCSBuildingHandler(userRoleService *service.FCSBuildingService) *FCSBuildingHandler {
	return &FCSBuildingHandler{userRoleService: userRoleService}
}

// GetFCSBuildings godoc
// @Summary Get a list of FCS Buildings
// @Description Retrieves FCS Buildings with pagination and optional search
// @Tags FCS Building
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering FCS Building"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/user-role [get]
func (h *FCSBuildingHandler) GetFCSBuildings(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search")
	offset := (page - 1) * rows

	total, err := h.userRoleService.GetTotal(search)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	fcsBuilding, err := h.userRoleService.GetAll(offset, rows, search)
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
		"items": fcsBuilding,
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

	if len(fcsBuilding) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetFCSBuilding godoc
// @Summary Get FCS Building by ID
// @Description Retrieve a specific FCS Building by its ID
// @Tags FCS Building
// @Accept json
// @Produce json
// @Param id path int true "FCS Building ID"
// @Success 200 {object} map[string]interface{} "FCS Building found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: FCS Building not found"
// @Router /admin/master/user-role/{id} [get]
func (h *FCSBuildingHandler) GetFCSBuilding(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	userRole, err := h.userRoleService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "FCS Building not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "FCS Building found successfully", userRole)
}

// CreateFCSBuilding godoc
// @Summary Create a new FCS Building
// @Description Create a new FCS Building with the provided details
// @Tags FCS Building
// @Accept json
// @Produce json
// @Param FCSBuilding body dto.FCSBuildings true "FCS Building details"
// @Success 201 {object} map[string]interface{} "FCS Building created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/user-role [post]
func (h *FCSBuildingHandler) UpdateFCSBuilding(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var input dto.FCSBuildings
	if err := c.BodyParser(&input); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var fcsBuilding []model.MstFCSBuilding
	for _, buildingID := range input.IDBuilding {
		fcsBuilding = append(fcsBuilding, model.MstFCSBuilding{
			IDFCS:      uint(ID),
			IDBuilding: buildingID,
		})
	}

	err = h.userRoleService.Delete(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	err = h.userRoleService.Create(&fcsBuilding)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusCreated, "FCS Building created successfully", fcsBuilding)
}
