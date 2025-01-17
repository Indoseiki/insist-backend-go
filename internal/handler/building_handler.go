package handler

import (
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type BuildingHandler struct {
	buildingService *service.BuildingService
}

func NewBuildingHandler(buildingService *service.BuildingService) *BuildingHandler {
	return &BuildingHandler{buildingService: buildingService}
}

// GetBuildings godoc
// @Summary Get a list of buildings
// @Description Retrieves buildings with pagination and optional search
// @Tags Building
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering building"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/building [get]
func (h *BuildingHandler) GetBuildings(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	idFCS := c.QueryInt("idFCS")
	plant := c.Query("plant")
	search := c.Query("search")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	total, err := h.buildingService.GetTotal(search, uint(idFCS), plant)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	buildings, err := h.buildingService.GetAll(offset, rows, search, sortBy, sortDirection, uint(idFCS), plant)
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
		"items": buildings,
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

	if len(buildings) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetBuilding godoc
// @Summary Get building by ID
// @Description Retrieve a specific building by its ID
// @Tags Building
// @Accept json
// @Produce json
// @Param id path int true "Building ID"
// @Success 200 {object} map[string]interface{} "Building found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Building not found"
// @Router /admin/master/building/{id} [get]
func (h *BuildingHandler) GetBuilding(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	building, err := h.buildingService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Building not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Building found successfully", building)
}

// CreateBuilding godoc
// @Summary Create a new building
// @Description Create a new building with the provided details
// @Tags Building
// @Accept json
// @Produce json
// @Param building body model.MstBuilding true "Building details"
// @Success 201 {object} map[string]interface{} "Building created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/building [post]
func (h *BuildingHandler) CreateBuilding(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var building model.MstBuilding
	if err := c.BodyParser(&building); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	building.IDCreatedby = userID
	building.IDUpdatedby = userID

	err := h.buildingService.Create(&building)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": building.ID,
	}

	return pkg.Response(c, fiber.StatusCreated, "Building created successfully", result)
}

// UpdateBuilding godoc
// @Summary Update an existing building
// @Description Update the details of an existing building by its ID
// @Tags Building
// @Accept json
// @Produce json
// @Param id path int true "Building ID"
// @Param building body model.MstBuilding true "Updated building details"
// @Success 200 {object} map[string]interface{} "Building updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 404 {object} map[string]interface{} "Not Found: Building not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/building/{id} [put]
func (h *BuildingHandler) UpdateBuilding(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var building *model.MstBuilding
	building, err = h.buildingService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Building not found"))
	}

	if err := c.BodyParser(building); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	building.ID = uint(ID)
	building.IDUpdatedby = userID

	err = h.buildingService.Update(building)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": building.ID,
	}

	return pkg.Response(c, fiber.StatusOK, "Building updated successfully", result)
}

// DeleteBuilding godoc
// @Summary Delete a building
// @Description Delete a building by its ID
// @Tags Building
// @Param id path int true "Building ID"
// @Success 200 {object} map[string]interface{} "Building deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Building not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/building/{id} [delete]
func (h *BuildingHandler) DeleteBuilding(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	user, err := h.buildingService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Building not found"))
	}

	err = h.buildingService.Delete(user)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Building deleted successfully", nil)
}
