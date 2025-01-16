package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UoMService struct {
	db *gorm.DB
}

func NewUoMService(db *gorm.DB) *UoMService {
	return &UoMService{db: db}
}

func (s *UoMService) GetByID(uomID uint) (*model.MstUoMs, error) {
	var uom model.MstUoMs
	if err := s.db.Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&uom, uomID).Error; err != nil {
		return nil, err
	}

	return &uom, nil
}

func (s *UoMService) GetTotal(search string) (int64, error) {
	var count int64

	query := s.db.Model(&model.MstUoMs{})

	if search != "" {
		query = query.Where("code ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *UoMService) GetAll(offset, limit int, search string, sortBy string, sortDirection bool) ([]model.MstUoMs, error) {
	var uoms []model.MstUoMs

	query := s.db.Model(&model.MstUoMs{}).Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
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

	if err := query.Find(&uoms).Error; err != nil {
		return nil, err
	}

	return uoms, nil
}

func (s *UoMService) Create(uom *model.MstUoMs) error {
	return s.db.Create(uom).Error
}

func (s *UoMService) Update(uom *model.MstUoMs) error {
	return s.db.Save(uom).Error
}

func (s *UoMService) Delete(uom *model.MstUoMs) error {
	return s.db.Delete(uom).Error
}
