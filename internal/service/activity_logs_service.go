package service

import (
	"insist-backend-golang/internal/model"
	"time"

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

func (s *ActivityLogService) GetTotal(search string, action string, isSuccess string, arrayDate []string) (int64, error) {
	var count int64

	query := s.db.Model(&model.ActivityLog{})

	if search != "" {
		query = query.Joins("JOIN mst_users u ON u.id = activity_logs.id_user").Where("u.name ILIKE ?", "%"+search+"%")
	}

	if action != "" {
		query = query.Where("action = ?", action)
	}

	if isSuccess != "" {
		query = query.Where("is_success =?", isSuccess)
	}

	if len(arrayDate) == 2 {
		startDate, err := time.Parse("2006-01-02", arrayDate[0])
		if err != nil {
			startDate = time.Time{}
		}

		endDate, err := time.Parse("2006-01-02", arrayDate[1])
		if err != nil {
			endDate = time.Time{}
		}

		if !startDate.IsZero() && !endDate.IsZero() {
			query = query.Where("DATE(created_at) BETWEEN ? AND ?", startDate, endDate)
		}
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *ActivityLogService) GetAll(offset, limit int, search string, action string, isSuccess string, sortBy string, sortDirection bool, arrayDate []string) ([]model.ActivityLog, error) {
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
		query = query.Joins("JOIN mst_users u ON u.id = activity_logs.id_user").Where("u.name ILIKE ?", "%"+search+"%")
	}

	if action != "" {
		query = query.Where("action = ?", action)
	}

	if isSuccess != "" {
		query = query.Where("is_success = ?", isSuccess)
	}

	if len(arrayDate) == 2 {
		startDate, err := time.Parse("2006-01-02", arrayDate[0])
		if err != nil {
			startDate = time.Time{}
		}

		endDate, err := time.Parse("2006-01-02", arrayDate[1])
		if err != nil {
			endDate = time.Time{}
		}

		if !startDate.IsZero() && !endDate.IsZero() {
			query = query.Where("DATE(created_at) BETWEEN ? AND ?", startDate, endDate)
		}
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

func (s *ActivityLogService) GetByUsername(username string) (*model.MstUser, error) {
	var user model.MstUser
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
