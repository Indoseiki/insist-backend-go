package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ItemService struct {
	db *gorm.DB
}

func NewItemService(db *gorm.DB) *ItemService {
	return &ItemService{db: db}
}

func (s *ItemService) GetByID(id uint) (*model.MstItem, error) {
	var item model.MstItem
	if err := s.db.Preload("ItemCategory").
		Preload("UOM").
		Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *ItemService) GetTotal(search string, idItemCategory uint, categoryCode string) (int64, error) {
	var count int64
	query := s.db.Model(&model.MstItem{}).Joins("LEFT JOIN mst_item_categories ON mst_item_categories.id = mst_items.id_item_category")

	if search != "" {
		query = query.Where("mst_items.code ILIKE ? OR mst_items.description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if idItemCategory != 0 {
		query = query.Where("mst_items.id_item_category = ?", idItemCategory)
	}
	if categoryCode != "" {
		query = query.Where("mst_item_categories.code = ?", categoryCode)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (s *ItemService) GetAll(offset, limit int, search, sortBy string, sortAsc bool, idItemCategory uint, categoryCode string) ([]model.MstItem, error) {
	var items []model.MstItem

	query := s.db.Model(&model.MstItem{}).
		Joins("LEFT JOIN mst_item_categories ON mst_item_categories.id = mst_items.id_item_category").
		Preload("ItemCategory").
		Preload("UOM").
		Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Offset(offset).Limit(limit)

	if sortBy != "" {
		query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: sortBy}, Desc: !sortAsc})
	}

	if search != "" {
		query = query.Where("mst_items.code ILIKE ? OR mst_items.description ILIKE ?", "%"+search+"%", "%"+search+"%")
	} else {
		query = query.Order("mst_items.code ASC")
	}

	if idItemCategory != 0 {
		query = query.Where("mst_items.id_item_category = ?", idItemCategory)
	}

	if categoryCode != "" {
		query = query.Where("mst_item_categories.code = ?", categoryCode)
	}

	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (s *ItemService) Create(item *model.MstItem) error {
	return s.db.Create(item).Error
}

func (s *ItemService) Update(item *model.MstItem) error {
	return s.db.Save(item).Error
}

func (s *ItemService) Delete(item *model.MstItem) error {
	return s.db.Delete(item).Error
}
