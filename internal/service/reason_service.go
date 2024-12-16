package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
)

type ReasonService struct {
	db *gorm.DB
}

func NewReasonService(db *gorm.DB) *ReasonService {
	return &ReasonService{db: db}
}

func (s *ReasonService) GetByID(reasonID uint) (*model.MstReason, error) {
	var reason model.MstReason
	if err := s.db.Preload("Menu", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, label")
	}).Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&reason, reasonID).Error; err != nil {
		return nil, err
	}

	return &reason, nil
}

func (s *ReasonService) GetTotal(search string, menuID uint) (int64, error) {
	var count int64

	query := s.db.Model(&model.MstReason{})

	if menuID != 0 {
		query = query.Where("id_menu = ?", menuID)
	}

	if search != "" {
		query = query.Where("key ILIKE ? OR code ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *ReasonService) GetAll(offset, limit int, search string, menuID uint) ([]model.MstReason, error) {
	var reasons []model.MstReason

	query := s.db.Model(&model.MstReason{}).Preload("Menu", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, label")
	}).Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Offset(offset).Limit(limit)

	if menuID != 0 {
		query = query.Where("id_menu = ?", menuID)
	}

	if search != "" {
		query = query.Where("key ILIKE ? OR code ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Find(&reasons).Error; err != nil {
		return nil, err
	}

	return reasons, nil
}

func (s *ReasonService) Create(reason *model.MstReason) error {
	return s.db.Create(reason).Error
}

func (s *ReasonService) Update(reason *model.MstReason) error {
	return s.db.Save(reason).Error
}

func (s *ReasonService) Delete(reason *model.MstReason) error {
	return s.db.Delete(reason).Error
}
