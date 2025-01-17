package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SectionService struct {
	db *gorm.DB
}

func NewSectionService(db *gorm.DB) *SectionService {
	return &SectionService{db: db}
}

func (s *SectionService) GetByID(buildingID uint) (*model.MstSection, error) {
	var building model.MstSection
	if err := s.db.Preload("FCS", func(db *gorm.DB) *gorm.DB {
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

func (s *SectionService) GetTotal(search string, IDFCS uint) (int64, error) {
	var count int64

	query := s.db.Model(&model.MstSection{})

	if IDFCS != 0 {
		query = query.Where("id_fcs = ?", IDFCS)
	}

	if search != "" {
		query = query.Where("code ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *SectionService) GetAll(offset, limit int, search string, sortBy string, sortDirection bool, IDFCS uint) ([]model.MstSection, error) {
	var buildings []model.MstSection

	query := s.db.Model(&model.MstSection{}).Preload("FCS", func(db *gorm.DB) *gorm.DB {
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

	if IDFCS != 0 {
		query = query.Where("id_fcs = ?", IDFCS)
	}

	if search != "" {
		query = query.Where("code ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Find(&buildings).Error; err != nil {
		return nil, err
	}

	return buildings, nil
}

func (s *SectionService) Create(building *model.MstSection) error {
	return s.db.Create(building).Error
}

func (s *SectionService) Update(building *model.MstSection) error {
	return s.db.Save(building).Error
}

func (s *SectionService) Delete(building *model.MstSection) error {
	return s.db.Delete(building).Error
}
