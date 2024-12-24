package handler

import (
	"fmt"
	"insist-backend-golang/internal/dto"
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type ApprovalUserHandler struct {
	approvalUserService *service.ApprovalUserService
}

func NewApprovalUserHandler(approvalUserService *service.ApprovalUserService) *ApprovalUserHandler {
	return &ApprovalUserHandler{approvalUserService: approvalUserService}
}

// GetApprovalUsers godoc
// @Summary Get a list of approval userss
// @Description Retrieves approval userss with pagination and optional search
// @Tags Approval User
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param rows query int false "Number of rows per page" default(20)
// @Param search query string false "Search keyword for filtering approval users"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/approval-user [get]
func (h *ApprovalUserHandler) GetApprovalUsers(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search")
	offset := (page - 1) * rows

	fmt.Println(page)

	total, err := h.approvalUserService.GetTotal(search)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	approvalUsers, err := h.approvalUserService.GetAll(offset, rows, search)
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
		"items": approvalUsers,
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

	if len(approvalUsers) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// CreateApprovalUser godoc
// @Summary Create a new approval users
// @Description Create a new approval users with the provided details
// @Tags Approval User
// @Accept json
// @Produce json
// @Param approvalUser body model.MstApprovalUser true "Approval User details"
// @Success 201 {object} map[string]interface{} "Approval User created successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/approval-user [post]
func (h *ApprovalUserHandler) UpdateApprovalUser(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	var input dto.ApprovalUsers
	if err := c.BodyParser(&input); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	if len(input.IDUser) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, "IDUser field cannot be empty"))
	}

	approvalUsers := make([]*model.MstApprovalUser, 0, len(input.IDUser))
	for _, userID := range input.IDUser {
		approvalUsers = append(approvalUsers, &model.MstApprovalUser{
			IDApproval: uint(ID),
			IDUser:     userID,
		})
	}

	if err := h.approvalUserService.Delete(uint(ID)); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	if err := h.approvalUserService.Create(approvalUsers); err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Approval User updated successfully", nil)
}

// DeleteApprovalUser godoc
// @Summary Delete a approval users
// @Description Delete a approval users by its ID
// @Tags Approval User
// @Param id path int true "Approval User ID"
// @Success 200 {object} map[string]interface{} "Approval User deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid ID"
// @Failure 404 {object} map[string]interface{} "Not Found: ApprovalUser not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/approval-user/{id} [delete]
func (h *ApprovalUserHandler) DeleteApprovalUser(c *fiber.Ctx) error {
	ID, err := c.ParamsInt("id")
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusBadRequest, err.Error()))
	}

	err = h.approvalUserService.Delete(uint(ID))
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	return pkg.Response(c, fiber.StatusOK, "Approval User deleted successfully", nil)
}
