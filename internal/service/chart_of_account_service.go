package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ChartOfAccountService struct {
	db *gorm.DB
}

func NewChartOfAccountService(db *gorm.DB) *ChartOfAccountService {
	return &ChartOfAccountService{db: db}
}

func (s *ChartOfAccountService) GetByID(deptID uint) (*model.MstChartOfAccount, error) {
	var dept model.MstChartOfAccount
	if err := s.db.Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&dept, deptID).Error; err != nil {
		return nil, err
	}

	return &dept, nil
}

func (s *ChartOfAccountService) GetTotal(search string) (int64, error) {
	var count int64

	query := s.db.Model(&model.MstChartOfAccount{})

	if search != "" {
		query = query.Where("account ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *ChartOfAccountService) GetAll(offset, limit int, search string, sortBy string, sortDirection bool) ([]model.MstChartOfAccount, error) {
	var depts []model.MstChartOfAccount

	query := s.db.Model(&model.MstChartOfAccount{}).Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
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
		query = query.Where("account ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Find(&depts).Error; err != nil {
		return nil, err
	}

	return depts, nil
}

func (s *ChartOfAccountService) Create(dept *model.MstChartOfAccount) error {
	return s.db.Create(dept).Error
}

func (s *ChartOfAccountService) Update(dept *model.MstChartOfAccount) error {
	return s.db.Save(dept).Error
}

func (s *ChartOfAccountService) Delete(dept *model.MstChartOfAccount) error {
	return s.db.Delete(dept).Error
}
