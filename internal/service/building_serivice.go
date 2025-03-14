package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BuildingService struct {
	db *gorm.DB
}

func NewBuildingService(db *gorm.DB) *BuildingService {
	return &BuildingService{db: db}
}

func (s *BuildingService) GetByID(buildingID uint) (*model.MstBuilding, error) {
	var building model.MstBuilding
	if err := s.db.Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&building, buildingID).Error; err != nil {
		return nil, err
	}

	return &building, nil
}

func (s *BuildingService) GetTotal(search string, IDFCS uint, plant string) (int64, error) {
	var count int64

	query := s.db.Model(&model.MstBuilding{}).Select("COUNT(DISTINCT mst_buildings.id)").Joins("LEFT JOIN mst_fcs_buildings ON mst_fcs_buildings.id_building = mst_buildings.id")

	if IDFCS != 0 {
		query = query.Where("mst_fcs_buildings.id_fcs = ?", IDFCS)
	}

	if plant != "" {
		query = query.Where("plant = ?", plant)
	}

	if search != "" {
		query = query.Where("code ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *BuildingService) GetAll(offset, limit int, search string, sortBy string, sortDirection bool, IDFCS uint, plant string) ([]model.MstBuilding, error) {
	var buildings []model.MstBuilding

	query := s.db.Model(&model.MstBuilding{}).Select("DISTINCT mst_buildings.*").Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Joins("LEFT JOIN mst_fcs_buildings ON mst_fcs_buildings.id_building = mst_buildings.id").Offset(offset).Limit(limit)

	if sortBy != "" {
		query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: sortBy}, Desc: sortDirection})
	} else {
		query = query.Order("updated_at ASC")
	}

	if IDFCS != 0 {
		query = query.Where("mst_fcs_buildings.id_fcs = ?", IDFCS)
	}

	if plant != "" {
		query = query.Where("plant = ?", plant)
	}

	if search != "" {
		query = query.Where("code ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Find(&buildings).Error; err != nil {
		return nil, err
	}

	return buildings, nil
}

func (s *BuildingService) Create(building *model.MstBuilding) error {
	return s.db.Create(building).Error
}

func (s *BuildingService) Update(building *model.MstBuilding) error {
	return s.db.Save(building).Error
}

func (s *BuildingService) Delete(building *model.MstBuilding) error {
	return s.db.Delete(building).Error
}
