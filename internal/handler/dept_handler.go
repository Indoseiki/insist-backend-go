package handler

import (
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type DeptHandler struct {
	deptService *service.DeptService
}

func NewDeptHandler(deptService *service.DeptService) *DeptHandler {
	return &DeptHandler{deptService: deptService}
}

// GetDepts godoc
// @Summary Get a list of departments
// @Description Retrieves departments with pagination and optional search
// @Tags Department
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering department"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/department [get]
func (h *DeptHandler) GetDepts(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	total, err := h.deptService.GetTotal(search)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	depts, err := h.deptService.GetAll(offset, rows, search, sortBy, sortDirection)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	totalPages := int(math.Ceil(float64(total) / float64(rows)))
	var nextPage *int
	if page < totalPages {
		nextPageVal := page + 1
		nextPage = &nextPageVal
	}

	result := map[string]interface{}{
		"items": depts,
		"pagination": map[string]interface{}{
			"current_page":  page,
			"next_page":     nextPage,
			"total_pages":   totalPages,
			"rows_per_page": rows,
			"total_rows":    total,
		},
	}

	if len(depts) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetDept godoc
// @Summary Get department by ID
// @Description Retrieve a specific department by its ID
// @Tags Department
// @Accept json
// @Produce json
// @Param id path int true "Department ID"
// @Success 200 {object} map[string]interface{} "Department found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Department not found"
// @Router /admin/master/department/{id} [get]
func (h *DeptHandler) GetDept(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	dept, err := h.deptService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Department not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Department found successfully", dept)
}

// CreateDept godoc
// @Summary Create a new department
// @Description Create a new department with the provided details
// @Tags Department
// @Accept json
// @Produce json
// @Param dept body model.MstDept true "Department details"
// @Success 201 {object} map[string]interface{} "Department created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/department [post]
func (h *DeptHandler) CreateDept(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var dept model.MstDept
	if err := c.BodyParser(&dept); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	dept.IDCreatedby = userID
	dept.IDUpdatedby = userID

	err := h.deptService.Create(&dept)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": dept.ID,
	}

	return pkg.Response(c, fiber.StatusCreated, "Department created successfully", result)
}

// UpdateDept godoc
// @Summary Update an existing department
// @Description Update the details of an existing department by its ID
// @Tags Department
// @Accept json
// @Produce json
// @Param id path int true "Department ID"
// @Param dept body model.MstDept true "Updated department details"
// @Success 200 {object} map[string]interface{} "Department updated successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 404 {object} map[string]interface{} "Not Found: Department not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/department/{id} [put]
func (h *DeptHandler) UpdateDept(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var dept *model.MstDept
	dept, err = h.deptService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Department not found"))
	}

	if err := c.BodyParser(dept); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	dept.ID = uint(ID)
	dept.IDUpdatedby = userID

	err = h.deptService.Update(dept)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	result := map[string]interface{}{
		"id": dept.ID,
	}

	return pkg.Response(c, fiber.StatusOK, "Department updated successfully", result)
}

// DeleteDept godoc
// @Summary Delete a department
// @Description Delete a department by its ID
// @Tags Department
// @Param id path int true "Department ID"
// @Success 200 {object} map[string]interface{} "Department deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: Department not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/department/{id} [delete]
func (h *DeptHandler) DeleteDept(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	user, err := h.deptService.GetByID(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Department not found"))
	}

	err = h.deptService.Delete(user)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Department deleted successfully", nil)
}
