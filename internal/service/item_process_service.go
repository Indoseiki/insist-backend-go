package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ItemProcessService struct {
	db *gorm.DB
}

func NewItemProcessService(db *gorm.DB) *ItemProcessService {
	return &ItemProcessService{db: db}
}

func (s *ItemProcessService) GetByID(id uint) (*model.MstItemProcess, error) {
	var itemProcess model.MstItemProcess
	if err := s.db.Preload("ItemCategory", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, code, description")
	}).Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&itemProcess, id).Error; err != nil {
		return nil, err
	}

	return &itemProcess, nil
}

func (s *ItemProcessService) GetTotal(search, categoryCode string) (int64, error) {
	var count int64

	query := s.db.Model(&model.MstItemProcess{}).
		Joins("LEFT JOIN mst_item_categories ON mst_item_categories.id = mst_item_processes.id_item_category")

	if categoryCode != "" {
		query = query.Where("mst_item_categories.code = ?", categoryCode)
	}

	if search != "" {
		query = query.Where("mst_item_processes.code ILIKE ? OR mst_item_processes.description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *ItemProcessService) GetAll(offset, limit int, search, sortBy string, sortDirection bool, categoryCode string) ([]model.MstItemProcess, error) {
	var itemProcesses []model.MstItemProcess

	query := s.db.Model(&model.MstItemProcess{}).
		Joins("LEFT JOIN mst_item_categories ON mst_item_categories.id = mst_item_processes.id_item_category").
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
		query = query.Order("mst_item_processes.code ASC")
	}

	if categoryCode != "" {
		query = query.Where("mst_item_categories.code = ?", categoryCode)
	}

	if search != "" {
		query = query.Where("mst_item_processes.code ILIKE ? OR mst_item_processes.description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Find(&itemProcesses).Error; err != nil {
		return nil, err
	}

	return itemProcesses, nil
}

func (s *ItemProcessService) Create(itemProcess *model.MstItemProcess) error {
	return s.db.Create(itemProcess).Error
}

func (s *ItemProcessService) Update(itemProcess *model.MstItemProcess) error {
	return s.db.Save(itemProcess).Error
}

func (s *ItemProcessService) Delete(itemProcess *model.MstItemProcess) error {
	return s.db.Delete(itemProcess).Error
}
