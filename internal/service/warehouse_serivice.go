package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WarehouseService struct {
	db *gorm.DB
}

func NewWarehouseService(db *gorm.DB) *WarehouseService {
	return &WarehouseService{db: db}
}

func (s *WarehouseService) GetByID(warehouseID uint) (*model.MstWarehouse, error) {
	var warehouse model.MstWarehouse
	if err := s.db.Preload("Building", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, code, description")
	}).Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&warehouse, warehouseID).Error; err != nil {
		return nil, err
	}

	return &warehouse, nil
}

func (s *WarehouseService) GetTotal(search string) (int64, error) {
	var count int64

	query := s.db.Model(&model.MstWarehouse{})

	if search != "" {
		query = query.Where("code ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *WarehouseService) GetAll(offset, limit int, search string, sortBy string, sortDirection bool) ([]model.MstWarehouse, error) {
	var warehouses []model.MstWarehouse

	query := s.db.Model(&model.MstWarehouse{}).Preload("Building", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, code, description")
	}).Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
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

	if err := query.Find(&warehouses).Error; err != nil {
		return nil, err
	}

	return warehouses, nil
}

func (s *WarehouseService) Create(warehouse *model.MstWarehouse) error {
	return s.db.Create(warehouse).Error
}

func (s *WarehouseService) Update(warehouse *model.MstWarehouse) error {
	return s.db.Save(warehouse).Error
}

func (s *WarehouseService) Delete(warehouse *model.MstWarehouse) error {
	return s.db.Delete(warehouse).Error
}
