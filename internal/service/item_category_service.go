package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ItemCategoryService struct {
	db *gorm.DB
}

func NewItemCategoryService(db *gorm.DB) *ItemCategoryService {
	return &ItemCategoryService{db: db}
}

func (s *ItemCategoryService) GetByID(bankID uint) (*model.MstItemCategory, error) {
	var bank model.MstItemCategory
	if err := s.db.Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&bank, bankID).Error; err != nil {
		return nil, err
	}

	return &bank, nil
}

func (s *ItemCategoryService) GetTotal(search string) (int64, error) {
	var count int64

	query := s.db.Model(&model.MstItemCategory{})

	if search != "" {
		query = query.Where("code ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *ItemCategoryService) GetAll(offset, limit int, search string, sortBy string, sortDirection bool) ([]model.MstItemCategory, error) {
	var banks []model.MstItemCategory

	query := s.db.Model(&model.MstItemCategory{}).Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
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

	if err := query.Find(&banks).Error; err != nil {
		return nil, err
	}

	return banks, nil
}

func (s *ItemCategoryService) Create(bank *model.MstItemCategory) error {
	return s.db.Create(bank).Error
}

func (s *ItemCategoryService) Update(bank *model.MstItemCategory) error {
	return s.db.Save(bank).Error
}

func (s *ItemCategoryService) Delete(bank *model.MstItemCategory) error {
	return s.db.Delete(bank).Error
}
