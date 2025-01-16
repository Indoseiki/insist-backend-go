package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProcessService struct {
	db *gorm.DB
}

func NewProcessService(db *gorm.DB) *ProcessService {
	return &ProcessService{db: db}
}

func (s *ProcessService) GetByID(processID uint) (*model.MstProcess, error) {
	var process model.MstProcess
	if err := s.db.Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&process, processID).Error; err != nil {
		return nil, err
	}

	return &process, nil
}

func (s *ProcessService) GetTotal(search string) (int64, error) {
	var count int64

	query := s.db.Model(&model.MstProcess{})

	if search != "" {
		query = query.Where("code ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *ProcessService) GetAll(offset, limit int, search string, sortBy string, sortDirection bool) ([]model.MstProcess, error) {
	var process []model.MstProcess

	query := s.db.Model(&model.MstProcess{}).Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
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

	if err := query.Find(&process).Error; err != nil {
		return nil, err
	}

	return process, nil
}

func (s *ProcessService) Create(process *model.MstProcess) error {
	return s.db.Create(process).Error
}

func (s *ProcessService) Update(process *model.MstProcess) error {
	return s.db.Save(process).Error
}

func (s *ProcessService) Delete(process *model.MstProcess) error {
	return s.db.Delete(process).Error
}
