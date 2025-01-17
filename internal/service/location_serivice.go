package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LocationService struct {
	db *gorm.DB
}

func NewLocationService(db *gorm.DB) *LocationService {
	return &LocationService{db: db}
}

func (s *LocationService) GetByID(locationID uint) (*model.MstLocation, error) {
	var location model.MstLocation
	if err := s.db.Preload("Warehouse", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, code, description")
	}).Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&location, locationID).Error; err != nil {
		return nil, err
	}

	return &location, nil
}

func (s *LocationService) GetTotal(search string, idWarehouse uint) (int64, error) {
	var count int64

	query := s.db.Model(&model.MstLocation{})

	if idWarehouse != 0 {
		query = query.Where("id_warehouse = ?", idWarehouse)
	}

	if search != "" {
		query = query.Where("location ILIKE ?", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *LocationService) GetAll(offset, limit int, search string, sortBy string, sortDirection bool, idWarehouse uint) ([]model.MstLocation, error) {
	var locations []model.MstLocation

	query := s.db.Model(&model.MstLocation{}).Preload("Warehouse", func(db *gorm.DB) *gorm.DB {
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

	if idWarehouse != 0 {
		query = query.Where("id_warehouse = ?", idWarehouse)
	}

	if search != "" {
		query = query.Where("location ILIKE ?", "%"+search+"%")
	}

	if err := query.Find(&locations).Error; err != nil {
		return nil, err
	}

	return locations, nil
}

func (s *LocationService) Create(location *model.MstLocation) error {
	return s.db.Create(location).Error
}

func (s *LocationService) Update(location *model.MstLocation) error {
	return s.db.Save(location).Error
}

func (s *LocationService) Delete(location *model.MstLocation) error {
	return s.db.Delete(location).Error
}
