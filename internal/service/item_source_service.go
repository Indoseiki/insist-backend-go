package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ItemSourceService struct {
	db *gorm.DB
}

func NewItemSourceService(db *gorm.DB) *ItemSourceService {
	return &ItemSourceService{db: db}
}

func (s *ItemSourceService) GetByID(id uint) (*model.MstItemSource, error) {
	var itemSource model.MstItemSource
	err := s.db.Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&itemSource, id).Error

	if err != nil {
		return nil, err
	}
	return &itemSource, nil
}

func (s *ItemSourceService) GetTotal(search, categoryCode string) (int64, error) {
	var count int64
	query := s.db.Model(&model.MstItemSource{}).
		Joins("LEFT JOIN mst_item_categories ON mst_item_categories.id = mst_item_sources.id_item_category")

	if categoryCode != "" {
		query = query.Where("mst_item_categories.code = ?", categoryCode)
	}

	if search != "" {
		query = query.Where("mst_item_sources.code ILIKE ? OR mst_item_sources.description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	err := query.Count(&count).Error
	return count, err
}

func (s *ItemSourceService) GetAll(offset, limit int, search, sortBy string, sortDirection bool, categoryCode string) ([]model.MstItemSource, error) {
	var itemSources []model.MstItemSource

	query := s.db.Model(&model.MstItemSource{}).
		Joins("LEFT JOIN mst_item_categories ON mst_item_categories.id = mst_item_sources.id_item_category").
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
		query = query.Order("mst_item_sources.code ASC")
	}

	if categoryCode != "" {
		query = query.Where("mst_item_categories.code = ?", categoryCode)
	}

	if search != "" {
		query = query.Where("mst_item_sources.code ILIKE ? OR mst_item_sources.description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	err := query.Find(&itemSources).Error
	return itemSources, err
}

func (s *ItemSourceService) Create(data *model.MstItemSource) error {
	return s.db.Create(data).Error
}

func (s *ItemSourceService) Update(data *model.MstItemSource) error {
	return s.db.Save(data).Error
}

func (s *ItemSourceService) Delete(data *model.MstItemSource) error {
	return s.db.Delete(data).Error
}
