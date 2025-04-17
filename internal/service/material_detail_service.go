package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MaterialDetailService struct {
	db *gorm.DB
}

func NewMaterialDetailService(db *gorm.DB) *MaterialDetailService {
	return &MaterialDetailService{db: db}
}

func (s *MaterialDetailService) GetByID(id uint) (*model.MstMaterialDetail, error) {
	var materialDetail model.MstMaterialDetail
	if err := s.db.Preload("Material").
		Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).
		Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).First(&materialDetail, id).Error; err != nil {
		return nil, err
	}
	return &materialDetail, nil
}

func (s *MaterialDetailService) GetTotal(search string, idMaterial uint) (int64, error) {
	var count int64
	query := s.db.Model(&model.MstMaterialDetail{}).
		Joins("LEFT JOIN mst_materials ON mst_materials.id = mst_material_details.id_material").
		Joins("LEFT JOIN mst_items ON mst_items.code = mst_materials.code")

	if search != "" {
		query = query.Where("mst_material_details.rmss_num ILIKE ? OR mst_items.code ILIKE ? OR mst_items.description ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if idMaterial != 0 {
		query = query.Where("mst_material_details.id_material = ?", idMaterial)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (s *MaterialDetailService) GetAll(offset, limit int, search, sortBy string, sortAsc bool, idMaterial uint) ([]model.MstMaterialDetail, error) {
	var materialDetails []model.MstMaterialDetail

	query := s.db.Model(&model.MstMaterialDetail{}).
		Joins("LEFT JOIN mst_materials ON mst_materials.id = mst_material_details.id_material").
		Joins("LEFT JOIN mst_items ON mst_items.code = mst_materials.code").
		Preload("Material").
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
		query = query.Where("mst_material_details.rmss_num ILIKE ? OR mst_items.code ILIKE ? OR mst_items.description ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	} else {
		query = query.Order("mst_items.code ASC")
	}

	if idMaterial != 0 {
		query = query.Where("mst_material_details.id_material = ?", idMaterial)
	}

	if err := query.Find(&materialDetails).Error; err != nil {
		return nil, err
	}
	return materialDetails, nil
}

func (s *MaterialDetailService) Create(materialDetail *model.MstMaterialDetail) error {
	return s.db.Create(materialDetail).Error
}

func (s *MaterialDetailService) Update(materialDetail *model.MstMaterialDetail) error {
	return s.db.Save(materialDetail).Error
}

func (s *MaterialDetailService) Delete(materialDetail *model.MstMaterialDetail) error {
	return s.db.Delete(materialDetail).Error
}
