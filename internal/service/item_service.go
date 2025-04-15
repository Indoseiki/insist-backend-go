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

func (s *ItemService) GetTotal(search, idItemCategory string) (int64, error) {
	var count int64
	query := s.db.Model(&model.MstItem{})

	if search != "" {
		query = query.Where("code ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if idItemCategory != "" {
		query = query.Where("id_item_category = ?", idItemCategory)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (s *ItemService) GetAll(offset, limit int, search, sortBy string, sortAsc bool, idItemCategory string) ([]model.MstItem, error) {
	var items []model.MstItem

	query := s.db.Model(&model.MstItem{}).
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
		query = query.Where("code ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if idItemCategory != "" {
		query = query.Where("id_item_category = ?", idItemCategory)
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
