package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FCSService struct {
	db *gorm.DB
}

func NewFCSService(db *gorm.DB) *FCSService {
	return &FCSService{db: db}
}

func (s *FCSService) GetByID(fcsID uint) (*model.MstFCS, error) {
	var fcs model.MstFCS
	if err := s.db.Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&fcs, fcsID).Error; err != nil {
		return nil, err
	}

	return &fcs, nil
}

func (s *FCSService) GetTotal(search string) (int64, error) {
	var count int64

	query := s.db.Model(&model.MstFCS{})

	if search != "" {
		query = query.Where("code ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *FCSService) GetAll(offset, limit int, search string, sortBy string, sortDirection bool) ([]model.MstFCS, error) {
	var fcs []model.MstFCS

	query := s.db.Model(&model.MstFCS{}).Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
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

	if err := query.Find(&fcs).Error; err != nil {
		return nil, err
	}

	return fcs, nil
}

func (s *FCSService) Create(fcs *model.MstFCS) error {
	return s.db.Create(fcs).Error
}

func (s *FCSService) Update(fcs *model.MstFCS) error {
	return s.db.Save(fcs).Error
}

func (s *FCSService) Delete(fcs *model.MstFCS) error {
	return s.db.Delete(fcs).Error
}
