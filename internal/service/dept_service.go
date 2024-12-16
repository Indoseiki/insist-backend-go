package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DeptService struct {
	db *gorm.DB
}

func NewDeptService(db *gorm.DB) *DeptService {
	return &DeptService{db: db}
}

func (s *DeptService) GetByID(deptID uint) (*model.MstDept, error) {
	var dept model.MstDept
	if err := s.db.Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&dept, deptID).Error; err != nil {
		return nil, err
	}

	return &dept, nil
}

func (s *DeptService) GetTotal(search string) (int64, error) {
	var count int64

	query := s.db.Model(&model.MstDept{})

	if search != "" {
		query = query.Where("code ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *DeptService) GetAll(offset, limit int, search string, sortBy string, sortDirection bool) ([]model.MstDept, error) {
	var depts []model.MstDept

	query := s.db.Model(&model.MstDept{}).Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Offset(offset).Limit(limit)

	if sortBy != "" {
		query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: sortBy}, Desc: sortDirection})
	} else {
		query = query.Order("updated_at ASC")
	}

	if search != "" {
		query = query.Where("code ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Find(&depts).Error; err != nil {
		return nil, err
	}

	return depts, nil
}

func (s *DeptService) Create(dept *model.MstDept) error {
	return s.db.Create(dept).Error
}

func (s *DeptService) Update(dept *model.MstDept) error {
	return s.db.Save(dept).Error
}

func (s *DeptService) Delete(dept *model.MstDept) error {
	return s.db.Delete(dept).Error
}
