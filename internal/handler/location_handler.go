package handler

import (
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type LocationHandler struct {
	locationService *service.LocationService
}

func NewLocationHandler(locationService *service.LocationService) *LocationHandler {
	return &LocationHandler{locationService: locationService}
}

// GetLocations godoc
// @Summary Get a list of locations
// @Description Retrieves locations with pagination and optional search
// @Tags Location
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering location"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /pid/master/location [get]
func (h *LocationHandler) GetLocations(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	idWarehouse := c.QueryInt("id_warehouse")
	search := c.Query("search")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	total, err := h.locationService.GetTotal(search, uint(idWarehouse))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	locations, err := h.locationService.GetAll(offset, rows, search, sortBy, sortDirection, uint(idWarehouse))
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
		"items": locations,
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

	if len(locations) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetLocation godoc
// @Summary Get location by ID
// @Description Retrieve a specific location by its ID
// @Tags Location
// @Accept json
// @Produce json
// @Param id path int true "Location ID"
// @Success 200 {object} map[string]interface{} "Location found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Location not found"
// @Router /pid/master/location/{id} [get]
func (h *LocationHandler) GetLocation(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	location, err := h.locationService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Location not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Location found successfully", location)
}

// CreateLocation godoc
// @Summary Create a new location
// @Description Create a new location with the provided details
// @Tags Location
// @Accept json
// @Produce json
// @Param location body model.MstLocation true "Location details"
// @Success 201 {object} map[string]interface{} "Location created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /pid/master/location [post]
func (h *LocationHandler) CreateLocation(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var location model.MstLocation
	if err := c.BodyParser(&location); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	location.IDCreatedby = userID
	location.IDUpdatedby = userID

	err := h.locationService.Create(&location)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": location.ID,
	}

	return pkg.Response(c, fiber.StatusCreated, "Location created successfully", result)
}

// UpdateLocation godoc
// @Summary Update an existing location
// @Description Update the details of an existing location by its ID
// @Tags Location
// @Accept json
// @Produce json
// @Param id path int true "Location ID"
// @Param location body model.MstLocation true "Updated location details"
// @Success 200 {object} map[string]interface{} "Location updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 404 {object} map[string]interface{} "Not Found: Location not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /pid/master/location/{id} [put]
func (h *LocationHandler) UpdateLocation(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var location *model.MstLocation
	location, err = h.locationService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Location not found"))
	}

	if err := c.BodyParser(location); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	location.ID = uint(ID)
	location.IDUpdatedby = userID

	err = h.locationService.Update(location)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": location.ID,
	}

	return pkg.Response(c, fiber.StatusOK, "Location updated successfully", result)
}

// DeleteLocation godoc
// @Summary Delete a location
// @Description Delete a location by its ID
// @Tags Location
// @Param id path int true "Location ID"
// @Success 200 {object} map[string]interface{} "Location deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Location not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /pid/master/location/{id} [delete]
func (h *LocationHandler) DeleteLocation(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	user, err := h.locationService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Location not found"))
	}

	err = h.locationService.Delete(user)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Location deleted successfully", nil)
}
