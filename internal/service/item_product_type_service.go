package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ItemProductTypeService struct {
	db *gorm.DB
}

func NewItemProductTypeService(db *gorm.DB) *ItemProductTypeService {
	return &ItemProductTypeService{db: db}
}

func (s *ItemProductTypeService) GetByID(itemProductTypeID uint) (*model.MstItemProductType, error) {
	var itemProductType model.MstItemProductType
	if err := s.db.Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&itemProductType, itemProductTypeID).Error; err != nil {
		return nil, err
	}

	return &itemProductType, nil
}

func (s *ItemProductTypeService) GetTotal(search string, idItemProduct uint) (int64, error) {
	var count int64

	query := s.db.Model(&model.MstItemProductType{})

	if idItemProduct != 0 {
		query = query.Where("id_item_product = ?", idItemProduct)
	}

	if search != "" {
		query = query.Where("code ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *ItemProductTypeService) GetAll(offset, limit int, search, sortBy string, sortDirection bool, idItemProduct uint) ([]model.MstItemProductType, error) {
	var itemProductTypes []model.MstItemProductType

	query := s.db.Model(&model.MstItemProductType{}).
		Preload("ItemProduct", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, code, description")
		}).
		Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).
		Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).
		Offset(offset).Limit(limit)

	if idItemProduct != 0 {
		query = query.Where("id_item_product = ?", idItemProduct)
	}

	if sortBy != "" {
		query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: sortBy}, Desc: sortDirection})
	} else {
		query = query.Order("code ASC")
	}

	if search != "" {
		query = query.Where("code ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Find(&itemProductTypes).Error; err != nil {
		return nil, err
	}

	return itemProductTypes, nil
}

func (s *ItemProductTypeService) Create(itemProductType *model.MstItemProductType) error {
	return s.db.Create(itemProductType).Error
}

func (s *ItemProductTypeService) Update(itemProductType *model.MstItemProductType) error {
	return s.db.Save(itemProductType).Error
}

func (s *ItemProductTypeService) Delete(itemProductType *model.MstItemProductType) error {
	return s.db.Delete(itemProductType).Error
}
