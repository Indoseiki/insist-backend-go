package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ItemSurfaceService struct {
	db *gorm.DB
}

func NewItemSurfaceService(db *gorm.DB) *ItemSurfaceService {
	return &ItemSurfaceService{db: db}
}

func (s *ItemSurfaceService) GetByID(id uint) (*model.MstItemSurface, error) {
	var itemSurface model.MstItemSurface
	if err := s.db.Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("ItemCategory", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, code, description")
	}).First(&itemSurface, id).Error; err != nil {
		return nil, err
	}
	return &itemSurface, nil
}

func (s *ItemSurfaceService) GetTotal(search, categoryCode string) (int64, error) {
	var count int64

	query := s.db.Model(&model.MstItemSurface{}).
		Joins("LEFT JOIN mst_item_categories ON mst_item_categories.id = mst_item_surfaces.id_item_category")

	if categoryCode != "" {
		query = query.Where("mst_item_categories.code = ?", categoryCode)
	}

	if search != "" {
		query = query.Where("mst_item_surfaces.code ILIKE ? OR mst_item_surfaces.description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *ItemSurfaceService) GetAll(offset, limit int, search, sortBy string, sortDirection bool, categoryCode string) ([]model.MstItemSurface, error) {
	var itemSurfaces []model.MstItemSurface

	query := s.db.Model(&model.MstItemSurface{}).
		Joins("LEFT JOIN mst_item_categories ON mst_item_categories.id = mst_item_surfaces.id_item_category").
		Preload("ItemCategory", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, code, description")
		}).Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Offset(offset).Limit(limit)

	if sortBy != "" {
		query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: sortBy}, Desc: sortDirection})
	} else {
		query = query.Order("mst_item_surfaces.code ASC")
	}

	if categoryCode != "" {
		query = query.Where("mst_item_categories.code = ?", categoryCode)
	}

	if search != "" {
		query = query.Where("mst_item_surfaces.code ILIKE ? OR mst_item_surfaces.description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Find(&itemSurfaces).Error; err != nil {
		return nil, err
	}

	return itemSurfaces, nil
}

func (s *ItemSurfaceService) Create(itemSurface *model.MstItemSurface) error {
	return s.db.Create(itemSurface).Error
}

func (s *ItemSurfaceService) Update(itemSurface *model.MstItemSurface) error {
	return s.db.Save(itemSurface).Error
}

func (s *ItemSurfaceService) Delete(itemSurface *model.MstItemSurface) error {
	return s.db.Delete(itemSurface).Error
}
