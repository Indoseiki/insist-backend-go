package handler

import (
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"

	"github.com/gofiber/fiber/v2"
)

type ItemInforHandler struct {
	itemInforService *service.ItemInforService
}

func NewItemInforHandler(itemInforService *service.ItemInforService) *ItemInforHandler {
	return &ItemInforHandler{itemInforService: itemInforService}
}

// GetItemInfors godoc
// @Summary Get list of Item Infors
// @Tags Item Infor
// @Accept json
// @Produce json
// @Param page query int false "Page"
// @Param rows query int false "Rows per page"
// @Param search query string false "Search"
// @Param sortBy query string false "Sort by field"
// @Param sortDirection query boolean false "true = ASC, false = DESC"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /infor/master/item [get]
func (h *ItemInforHandler) GetItemInfors(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search", "")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection", true)
	offset := (page - 1) * rows

	total, err := h.itemInforService.GetTotal(search)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	items, err := h.itemInforService.GetAll(offset, rows, search, sortBy, sortDirection)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	totalPages := int(math.Ceil(float64(total) / float64(rows)))
	start := offset + 1
	end := int(math.Min(float64(offset+rows), float64(total)))

	result := map[string]interface{}{
		"items": items,
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
