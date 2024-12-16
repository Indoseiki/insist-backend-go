package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EmployeeService struct {
	db *gorm.DB
}

func NewEmployeeService(db *gorm.DB) *EmployeeService {
	return &EmployeeService{db: db}
}

func (s *EmployeeService) GetByNumber(number string) (*model.MstEmployee, error) {
	var employee model.MstEmployee
	if err := s.db.First(&employee, number).Error; err != nil {
		return nil, err
	}

	return &employee, nil
}

func (s *EmployeeService) GetTotal(search string) (int64, error) {
	var count int64

	query := s.db.Model(&model.MstEmployee{})

	if search != "" {
		query = query.Where("number ILIKE ? OR name ILIKE ? OR division ILIKE ? OR department ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *EmployeeService) GetAll(offset, limit int, search string, sortBy string, sortDirection bool) ([]model.MstEmployee, error) {
	var depts []model.MstEmployee

	query := s.db.Model(&model.MstEmployee{}).Offset(offset).Limit(limit)

	if sortBy != "" {
		query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: sortBy}, Desc: sortDirection})
	} else {
		query = query.Order("updated_at ASC")
	}

	if search != "" {
		query = query.Where("number ILIKE ? OR name ILIKE ? OR division ILIKE ? OR department ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Find(&depts).Error; err != nil {
		return nil, err
	}

	return depts, nil
}

func (s *EmployeeService) Create(employee *model.MstEmployee) error {
	return s.db.Create(employee).Error
}

func (s *EmployeeService) Update(employee *model.MstEmployee) error {
	return s.db.Where("number = ?", employee.Number).Updates(employee).Error
}

func (s *EmployeeService) SetInactiveIfNotInList(employeeNumbers []string) error {
	return s.db.Model(&model.MstEmployee{}).
		Where("number NOT IN ?", employeeNumbers).
		Update("is_active", false).Error
}
