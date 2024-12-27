package handler

import (
	"encoding/json"
	"fmt"
	"insist-backend-golang/internal/model"
	"insist-backend-golang/internal/service"
	"insist-backend-golang/pkg"
	"math"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
)

type EmployeeHandler struct {
	employeeService *service.EmployeeService
}

func NewEmployeeHandler(employeeService *service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{employeeService: employeeService}
}

// GetEmployees godoc
// @Summary Retrieve Employee List
// @Description Fetch a paginated list of employees with optional search parameters
// @Tags Employee
// @Accept json
// @Produce json
// @Param page query int false "Page number (default is 1)"
// @Param rows query int false "Number of rows per page (default is 20)"
// @Param search query string false "Search term for employee name or number"
// @Success 200 {object} map[string]interface{} "Data found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid query parameters"
// @Failure 404 {object} map[string]interface{} "Not Found: No data found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/employee [get]
func (h *EmployeeHandler) GetEmployees(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	rows := c.QueryInt("rows", 20)
	search := c.Query("search", "")
	sortBy := c.Query("sortBy", "")
	sortDirection := c.QueryBool("sortDirection")
	offset := (page - 1) * rows

	total, err := h.employeeService.GetTotal(search)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	employees, err := h.employeeService.GetAll(offset, rows, search, sortBy, sortDirection)
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
		"items": employees,
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

	if len(employees) == 0 {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "No data found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Data found successfully", result)
}

// GetEmployee godoc
// @Summary Retrieve Employee by Number
// @Description Fetch the details of a specific employee using their employee number
// @Tags Employee
// @Accept json
// @Produce json
// @Param number path int true "Employee number"
// @Success 200 {object} map[string]interface{} "Employee found successfully"
// @Failure 400 {object} map[string]interface{} "Bad Request: Invalid employee number"
// @Failure 404 {object} map[string]interface{} "Not Found: Employee not found"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/employee/{number} [get]
func (h *EmployeeHandler) GetEmployee(c *fiber.Ctx) error {
	number := c.Params("number")

	fmt.Println(number)
	employee, err := h.employeeService.GetByNumber(number)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusNotFound, "Employee not found"))
	}

	return pkg.Response(c, fiber.StatusOK, "Employee found successfully", employee)
}

// SyncEmployee godoc
// @Summary Synchronize Employee Data
// @Description Sync employee data from the external HRIS API and update the local database
// @Tags Employee
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Data synchronization successful"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /admin/master/employee/sync [get]
func (h *EmployeeHandler) SyncEmployee(c *fiber.Ctx) error {
	baseUrl := "https://indoseiki.hris-server.com/api/employees/reads"
	alphabet := "abcdefghijklmnopqrstuvwxyz"

	cookies := []string{
		"wssplashuid=be71f706e274c47df22311b8f50e92026cc4a4e2.1714383404.1",
		"ci_session=8b2a26ab1896853ce33ca982c0a3d26b23c4ec7f",
	}

	uniqueEmployeeNumbers := make(map[string]bool)
	var uniqueEmployees []model.MstEmployee

	client := resty.New()

	for _, letter := range alphabet {
		url := fmt.Sprintf("%s?number=&name=%s", baseUrl, string(letter))

		resp, err := client.R().
			SetHeader("Cookie", fmt.Sprintf("%s", cookies)).
			Get(url)

		if err != nil {
			return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
		}

		var jsonData struct {
			Code    int                      `json:"code"`
			Results []map[string]interface{} `json:"results"`
		}

		if err := json.Unmarshal(resp.Body(), &jsonData); err != nil {
			return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
		}

		for _, result := range jsonData.Results {
			var employee model.MstEmployee
			number, ok := result["number"].(string)
			if !ok || number == "" {
				continue
			}
			if uniqueEmployeeNumbers[number] {
				continue
			}

			uniqueEmployeeNumbers[number] = true
			employee.Number = number

			employee.Name, _ = result["name"].(string)
			employee.Division, _ = result["division"].(string)
			employee.Department, _ = result["departement"].(string)
			employee.Position, _ = result["position"].(string)
			employee.Service, _ = result["service"].(string)
			employee.Education, _ = result["education"].(string)

			if birthdayStr, ok := result["birthday"].(string); ok && birthdayStr != "" {
				if birthdayStr != "0000-00-00" && birthdayStr != "1899-11-30" {
					birthday, err := time.Parse("2006-01-02", birthdayStr)
					if err == nil {
						employee.Birthday = birthday
					}
				}
			}

			employee.IsActive = true
			uniqueEmployees = append(uniqueEmployees, employee)
		}
	}

	for _, employee := range uniqueEmployees {
		nik := employee.Number

		existingEmployee, _ := h.employeeService.GetByNumber(nik)

		if existingEmployee != nil {
			err := h.employeeService.Update(&employee)
			if err != nil {
				return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
			}
		} else {
			err := h.employeeService.Create(&employee)
			if err != nil {
				return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
			}
		}

	}

	allEmployeeNumbers := make([]string, 0, len(uniqueEmployeeNumbers))
	for number := range uniqueEmployeeNumbers {
		allEmployeeNumbers = append(allEmployeeNumbers, number)
	}

	err := h.employeeService.SetInactiveIfNotInList(allEmployeeNumbers)
	if err != nil {
		return pkg.ErrorResponse(c, fiber.NewError(fiber.StatusInternalServerError, err.Error()))
	}

	fmt.Println("Employee synchronization successfully")
	return pkg.Response(c, fiber.StatusOK, "Employee synchronization successfully", nil)
}
