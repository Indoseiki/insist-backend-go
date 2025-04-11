package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ItemProductService struct {
	db *gorm.DB
}

func NewItemProductService(db *gorm.DB) *ItemProductService {
	return &ItemProductService{db: db}
}

func (s *ItemProductService) GetByID(itemProductID uint) (*model.MstItemProduct, error) {
	var itemProduct model.MstItemProduct
	if err := s.db.Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&itemProduct, itemProductID).Error; err != nil {
		return nil, err
	}

	return &itemProduct, nil
}

func (s *ItemProductService) GetTotal(search, categoryCode, subCategoryCode string) (int64, error) {
	var count int64

	query := s.db.Model(&model.MstItemProduct{}).
		Joins("LEFT JOIN mst_item_categories ON mst_item_categories.id = mst_item_products.id_item_category").
		Joins("LEFT JOIN mst_item_sub_categories ON mst_item_sub_categories.id = mst_item_products.id_item_sub_category")

	if categoryCode != "" {
		query = query.Where("mst_item_categories.code = ?", categoryCode)
	}

	if subCategoryCode != "" {
		query = query.Where("mst_item_sub_categories.code = ?", subCategoryCode)
	}

	if search != "" {
		query = query.Where("mst_item_products.code ILIKE ? OR mst_item_products.description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *ItemProductService) GetAll(offset, limit int, search, sortBy string, sortDirection bool, categoryCode, subCategoryCode string) ([]model.MstItemProduct, error) {
	var itemProducts []model.MstItemProduct

	query := s.db.Model(&model.MstItemProduct{}).
		Joins("LEFT JOIN mst_item_categories ON mst_item_categories.id = mst_item_products.id_item_category").
		Joins("LEFT JOIN mst_item_sub_categories ON mst_item_sub_categories.id = mst_item_products.id_item_sub_category").
		Preload("ItemCategory", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, code, description")
		}).Preload("ItemSubCategory", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, code, description")
	}).Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).
		Offset(offset).Limit(limit)

	if sortBy != "" {
		query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: sortBy}, Desc: sortDirection})
	} else {
		query = query.Order("mst_item_products.updated_at ASC")
	}

	if categoryCode != "" {
		query = query.Where("mst_item_categories.code = ?", categoryCode)
	}

	if subCategoryCode != "" {
		query = query.Where("mst_item_sub_categories.code = ?", subCategoryCode)
	}

	if search != "" {
		query = query.Where("mst_item_products.code ILIKE ? OR mst_item_products.description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Find(&itemProducts).Error; err != nil {
		return nil, err
	}

	return itemProducts, nil
}

func (s *ItemProductService) Create(itemProduct *model.MstItemProduct) error {
	return s.db.Create(itemProduct).Error
}

func (s *ItemProductService) Update(itemProduct *model.MstItemProduct) error {
	return s.db.Save(itemProduct).Error
}

func (s *ItemProductService) Delete(itemProduct *model.MstItemProduct) error {
	return s.db.Delete(itemProduct).Error
}
