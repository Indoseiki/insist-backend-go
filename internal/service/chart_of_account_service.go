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

func (s *ChartOfAccountService) GetByID(chartOfAccountID uint) (*model.MstChartOfAccount, error) {
	var chartOfAccount model.MstChartOfAccount
	if err := s.db.Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&chartOfAccount, chartOfAccountID).Error; err != nil {
		return nil, err
	}

	return &chartOfAccount, nil
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
	var chartOfAccounts []model.MstChartOfAccount

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

	if err := query.Find(&chartOfAccounts).Error; err != nil {
		return nil, err
	}

	return chartOfAccounts, nil
}

func (s *ChartOfAccountService) Create(chartOfAccount *model.MstChartOfAccount) error {
	return s.db.Create(chartOfAccount).Error
}

func (s *ChartOfAccountService) Update(chartOfAccount *model.MstChartOfAccount) error {
	return s.db.Save(chartOfAccount).Error
}

func (s *ChartOfAccountService) Delete(chartOfAccount *model.MstChartOfAccount) error {
	return s.db.Delete(chartOfAccount).Error
}
