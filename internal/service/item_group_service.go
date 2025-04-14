package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ItemGroupService struct {
	db *gorm.DB
}

func NewItemGroupService(db *gorm.DB) *ItemGroupService {
	return &ItemGroupService{db: db}
}

func (s *ItemGroupService) GetByID(itemGroupID uint) (*model.MstItemGroup, error) {
	var itemGroup model.MstItemGroup
	if err := s.db.Preload("ItemProductType", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, code, description")
	}).Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&itemGroup, itemGroupID).Error; err != nil {
		return nil, err
	}

	return &itemGroup, nil
}

func (s *ItemGroupService) GetTotal(search, categoryCode string, idProductType uint) (int64, error) {
	var count int64

	query := s.db.Model(&model.MstItemGroup{}).
		Joins("LEFT JOIN mst_item_product_types ON mst_item_product_types.id = mst_item_groups.id_item_product_type").
		Joins("LEFT JOIN mst_item_products ON mst_item_products.id = mst_item_product_types.id_item_product").
		Joins("LEFT JOIN mst_item_categories ON mst_item_categories.id = mst_item_products.id_item_category")

	if idProductType != 0 {
		query = query.Where("mst_item_product_types.id = ?", idProductType)
	}

	if categoryCode != "" {
		query = query.Where("mst_item_categories.code = ?", categoryCode)
	}

	if search != "" {
		query = query.Where("mst_item_groups.code ILIKE ? OR mst_item_groups.description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *ItemGroupService) GetAll(offset, limit int, search, sortBy string, sortDirection bool, categoryCode string, idProductType uint) ([]model.MstItemGroup, error) {
	var itemGroups []model.MstItemGroup

	query := s.db.Model(&model.MstItemGroup{}).
		Joins("LEFT JOIN mst_item_product_types ON mst_item_product_types.id = mst_item_groups.id_item_product_type").
		Joins("LEFT JOIN mst_item_products ON mst_item_products.id = mst_item_product_types.id_item_product").
		Joins("LEFT JOIN mst_item_categories ON mst_item_categories.id = mst_item_products.id_item_category").
		Preload("ItemProductType", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, code, description")
		}).Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Offset(offset).Limit(limit)

	if sortBy != "" {
		query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: sortBy}, Desc: sortDirection})
	} else {
		query = query.Order("mst_item_groups.updated_at ASC")
	}

	if idProductType != 0 {
		query = query.Where("mst_item_product_types.id = ?", idProductType)
	}

	if categoryCode != "" {
		query = query.Where("mst_item_categories.code = ?", categoryCode)
	}

	if search != "" {
		query = query.Where("mst_item_groups.code ILIKE ? OR mst_item_groups.description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Find(&itemGroups).Error; err != nil {
		return nil, err
	}

	return itemGroups, nil
}

func (s *ItemGroupService) Create(itemGroup *model.MstItemGroup) error {
	return s.db.Create(itemGroup).Error
}

func (s *ItemGroupService) Update(itemGroup *model.MstItemGroup) error {
	return s.db.Save(itemGroup).Error
}

func (s *ItemGroupService) Delete(itemGroup *model.MstItemGroup) error {
	return s.db.Delete(itemGroup).Error
}
