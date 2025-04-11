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

func (s *ItemCategoryService) GetByID(itemCategoryID uint) (*model.MstItemCategory, error) {
	var itemCategory model.MstItemCategory
	if err := s.db.Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&itemCategory, itemCategoryID).Error; err != nil {
		return nil, err
	}

	return &itemCategory, nil
}

func (s *ItemCategoryService) GetByCode(code string) (*model.MstItemCategory, error) {
	var itemCategory model.MstItemCategory
	if err := s.db.Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Where("code = ?", code).First(&itemCategory).Error; err != nil {
		return nil, err
	}

	return &itemCategory, nil
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
	var itemCategorys []model.MstItemCategory

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

	if err := query.Find(&itemCategorys).Error; err != nil {
		return nil, err
	}

	return itemCategorys, nil
}

func (s *ItemCategoryService) Create(itemCategory *model.MstItemCategory) error {
	return s.db.Create(itemCategory).Error
}

func (s *ItemCategoryService) Update(itemCategory *model.MstItemCategory) error {
	return s.db.Save(itemCategory).Error
}

func (s *ItemCategoryService) Delete(itemCategory *model.MstItemCategory) error {
	return s.db.Delete(itemCategory).Error
}
