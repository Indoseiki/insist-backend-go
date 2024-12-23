package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ActivityLogService struct {
	db *gorm.DB
}

func NewActivityLogService(db *gorm.DB) *ActivityLogService {
	return &ActivityLogService{db: db}
}

func (s *ActivityLogService) GetByID(logID uint) (*model.ActivityLog, error) {
	var log model.ActivityLog
	if err := s.db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&log, logID).Error; err != nil {
		return nil, err
	}

	return &log, nil
}

func (s *ActivityLogService) GetTotal(search string) (int64, error) {
	var count int64

	query := s.db.Model(&model.ActivityLog{})

	if search != "" {
		query = query.Where("action ILIKE ?", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *ActivityLogService) GetAll(offset, limit int, search string, sortBy string, sortDirection bool) ([]model.ActivityLog, error) {
	var logs []model.ActivityLog

	query := s.db.Model(&model.ActivityLog{}).Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Offset(offset).Limit(limit)

	if sortBy != "" {
		query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: sortBy}, Desc: sortDirection})
	} else {
		query = query.Order("created_at ASC")
	}

	if search != "" {
		query = query.Where("action ILIKE ?", "%"+search+"%")
	}

	if err := query.Find(&logs).Error; err != nil {
		return nil, err
	}

	return logs, nil
}

func (s *ActivityLogService) Create(log *model.ActivityLog) error {
	return s.db.Create(log).Error
}

func (s *ActivityLogService) Update(log *model.ActivityLog) error {
	return s.db.Save(log).Error
}

func (s *ActivityLogService) Delete(log *model.ActivityLog) error {
	return s.db.Delete(log).Error
}
