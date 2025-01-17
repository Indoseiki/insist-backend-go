package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SubSectionService struct {
	db *gorm.DB
}

func NewSubSectionService(db *gorm.DB) *SubSectionService {
	return &SubSectionService{db: db}
}

func (s *SubSectionService) GetByID(buildingID uint) (*model.MstSubSection, error) {
	var building model.MstSubSection
	if err := s.db.Preload("Section", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, id_fcs, code, description")
	}).Preload("Section.FCS", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, code, description")
	}).Preload("Building", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, code, description")
	}).Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&building, buildingID).Error; err != nil {
		return nil, err
	}

	return &building, nil
}

func (s *SubSectionService) GetTotal(search string, idSection uint) (int64, error) {
	var count int64

	query := s.db.Model(&model.MstSubSection{})

	if idSection != 0 {
		query = query.Where("id_section = ?", idSection)
	}

	if search != "" {
		query = query.Where("code ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *SubSectionService) GetAll(offset, limit int, search string, sortBy string, sortDirection bool, idSection uint) ([]model.MstSubSection, error) {
	var buildings []model.MstSubSection

	query := s.db.Model(&model.MstSubSection{}).Preload("Section", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, id_fcs, code, description")
	}).Preload("Section.FCS", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, code, description")
	}).Preload("Building", func(db *gorm.DB) *gorm.DB {
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

	if idSection != 0 {
		query = query.Where("id_section = ?", idSection)
	}

	if search != "" {
		query = query.Where("code ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Find(&buildings).Error; err != nil {
		return nil, err
	}

	return buildings, nil
}

func (s *SubSectionService) Create(building *model.MstSubSection) error {
	return s.db.Create(building).Error
}

func (s *SubSectionService) Update(building *model.MstSubSection) error {
	return s.db.Save(building).Error
}

func (s *SubSectionService) Delete(building *model.MstSubSection) error {
	return s.db.Delete(building).Error
}
