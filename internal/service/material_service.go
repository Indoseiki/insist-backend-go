package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MaterialService struct {
	db *gorm.DB
}

func NewMaterialService(db *gorm.DB) *MaterialService {
	return &MaterialService{db: db}
}

func (s *MaterialService) GetByID(id uint) (*model.MstMaterial, error) {
	var material model.MstMaterial
	if err := s.db.Preload("Item").
		Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).
		Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).First(&material, id).Error; err != nil {
		return nil, err
	}
	return &material, nil
}

func (s *MaterialService) GetTotal(search string) (int64, error) {
	var count int64
	query := s.db.Model(&model.MstMaterial{}).
		Joins("LEFT JOIN mst_items ON mst_items.code = mst_materials.code")

	if search != "" {
		query = query.Where("mst_items.code ILIKE ? OR mst_items.description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (s *MaterialService) GetAll(offset, limit int, search, sortBy string, sortAsc bool) ([]model.MstMaterial, error) {
	var materials []model.MstMaterial

	query := s.db.Model(&model.MstMaterial{}).
		Joins("LEFT JOIN mst_items ON mst_items.code = mst_materials.code").
		Preload("Item").
		Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).
		Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).
		Offset(offset).
		Limit(limit)

	if sortBy != "" {
		query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: sortBy}, Desc: !sortAsc})
	}

	if search != "" {
		query = query.Where("mst_items.code ILIKE ? OR mst_items.description ILIKE ?", "%"+search+"%", "%"+search+"%")
	} else {
		query = query.Order("mst_materials.code ASC")
	}

	if err := query.Find(&materials).Error; err != nil {
		return nil, err
	}
	return materials, nil
}

func (s *MaterialService) Create(material *model.MstMaterial) error {
	return s.db.Create(material).Error
}

func (s *MaterialService) Update(material *model.MstMaterial) error {
	return s.db.Save(material).Error
}

func (s *MaterialService) Delete(material *model.MstMaterial) error {
	return s.db.Delete(material).Error
}
