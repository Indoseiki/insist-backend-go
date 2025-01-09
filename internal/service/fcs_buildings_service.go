package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
)

type FCSBuildingService struct {
	db *gorm.DB
}

func NewFCSBuildingService(db *gorm.DB) *FCSBuildingService {
	return &FCSBuildingService{db: db}
}

func (s *FCSBuildingService) GetByID(fcsID uint) (*model.MstFCS, error) {
	var fcsBuilding model.MstFCS
	if err := s.db.Select("*").Preload("FCSBuilding", func(db *gorm.DB) *gorm.DB {
		return db.Select("id_fcs, id_building")
	}).Preload("FCSBuilding.Building", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, code, description, plant")
	}).Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&fcsBuilding, fcsID).Error; err != nil {
		return nil, err
	}

	return &fcsBuilding, nil
}

func (s *FCSBuildingService) GetTotal(search string) (int64, error) {
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

func (s *FCSBuildingService) GetAll(offset, limit int, search string) ([]model.MstFCS, error) {
	var fcsBuildings []model.MstFCS

	query := s.db.Model(&model.MstFCS{}).Preload("FCSBuilding", func(db *gorm.DB) *gorm.DB {
		return db.Select("id_fcs, id_building")
	}).Preload("FCSBuilding.Building", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, code, description, plant")
	}).Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Select("*").Offset(offset).Limit(limit)

	if search != "" {
		query = query.Where("code ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Find(&fcsBuildings).Error; err != nil {
		return nil, err
	}

	return fcsBuildings, nil
}

func (s *FCSBuildingService) Create(fcsBuilding *[]model.MstFCSBuilding) error {
	return s.db.Create(fcsBuilding).Error
}

func (s *FCSBuildingService) Update(fcsBuilding *model.MstFCSBuilding) error {
	return s.db.Save(fcsBuilding).Error
}

func (s *FCSBuildingService) Delete(fcsID uint) error {
	return s.db.Where("id_fcs = ?", fcsID).Delete(&model.MstFCSBuilding{}).Error
}
