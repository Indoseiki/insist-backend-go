package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ItemRawMaterialService struct {
	db *gorm.DB
}

func NewItemRawMaterialService(db *gorm.DB) *ItemRawMaterialService {
	return &ItemRawMaterialService{db: db}
}

func (s *ItemRawMaterialService) GetByID(id uint) (*model.MstItemRawMaterial, error) {
	var itemRawMaterial model.MstItemRawMaterial
	if err := s.db.Preload("Item").
		Preload("ItemProductType").
		Preload("ItemGroupType").
		Preload("ItemProcess").
		Preload("ItemSurface").
		Preload("ItemSource").
		Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&itemRawMaterial, id).Error; err != nil {
		return nil, err
	}
	return &itemRawMaterial, nil
}

func (s *ItemRawMaterialService) GetTotal(search, categoryCode string) (int64, error) {
	var count int64
	query := s.db.Model(&model.MstItemRawMaterial{}).
		Joins("JOIN mst_items ON mst_items.id = mst_item_raw_materials.id_item").
		Joins("JOIN mst_item_categories ON mst_items.id_item_category = mst_item_categories.id")

	if search != "" {
		query = query.Where("mst_items.code ILIKE ? OR mst_items.description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if categoryCode != "" {
		query = query.Where("mst_item_categories.code = ?", categoryCode)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (s *ItemRawMaterialService) GetAll(offset, limit int, search, categoryCode, sortBy string, sortAsc bool) ([]model.MstItemRawMaterial, error) {
	var itemRawMaterials []model.MstItemRawMaterial

	query := s.db.Model(&model.MstItemRawMaterial{}).
		Joins("JOIN mst_items ON mst_items.id = mst_item_raw_materials.id_item").
		Joins("JOIN mst_item_categories ON mst_items.id_item_category = mst_item_categories.id").
		Preload("Item").
		Preload("ItemProductType").
		Preload("ItemGroupType").
		Preload("ItemProcess").
		Preload("ItemSurface").
		Preload("ItemSource").
		Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Offset(offset).Limit(limit)

	if sortBy != "" {
		query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: sortBy}, Desc: !sortAsc})
	} else {
		query = query.Order("mst_items.code ASC")
	}

	if search != "" {
		query = query.Where("mst_items.code ILIKE ? OR mst_items.description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if categoryCode != "" {
		query = query.Where("mst_item_categories.code = ?", categoryCode)
	}

	if err := query.Find(&itemRawMaterials).Error; err != nil {
		return nil, err
	}
	return itemRawMaterials, nil
}

func (s *ItemRawMaterialService) Create(itemRawMaterial *model.MstItemRawMaterial) error {
	return s.db.Create(itemRawMaterial).Error
}

func (s *ItemRawMaterialService) Update(itemRawMaterial *model.MstItemRawMaterial) error {
	return s.db.Save(itemRawMaterial).Error
}

func (s *ItemRawMaterialService) Delete(itemRawMaterial *model.MstItemRawMaterial) error {
	return s.db.Delete(itemRawMaterial).Error
}
