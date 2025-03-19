package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ApprovalHistoryService struct {
	db *gorm.DB
}

func NewApprovalHistoryService(db *gorm.DB) *ApprovalHistoryService {
	return &ApprovalHistoryService{db: db}
}

func (s *ApprovalHistoryService) GetByID(approvalHistoryID uint) (*model.ApprovalHistory, error) {
	var approvalHistory model.ApprovalHistory
	if err := s.db.Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&approvalHistory, approvalHistoryID).Error; err != nil {
		return nil, err
	}

	return &approvalHistory, nil
}

func (s *ApprovalHistoryService) GetTotal(search string) (int64, error) {
	var count int64

	query := s.db.Model(&model.ApprovalHistory{})

	if search != "" {
		query = query.Where("key ILIKE ?", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *ApprovalHistoryService) GetAll(offset, limit int, search string, sortBy string, sortDirection bool) ([]model.ApprovalHistory, error) {
	var approvalHistories []model.ApprovalHistory

	query := s.db.Model(&model.ApprovalHistory{}).Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Offset(offset).Limit(limit)

	if search != "" {
		query = query.Where("key ILIKE ?", "%"+search+"%")
	}

	if sortBy != "" {
		query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: sortBy}, Desc: sortDirection})
	} else {
		query = query.Order("created_at ASC")
	}

	if err := query.Find(&approvalHistories).Error; err != nil {
		return nil, err
	}

	return approvalHistories, nil
}

func (s *ApprovalHistoryService) Create(approvalHistory *model.ApprovalHistory) error {
	return s.db.Create(approvalHistory).Error
}

func (s *ApprovalHistoryService) Update(approvalHistory *model.ApprovalHistory) error {
	return s.db.Save(approvalHistory).Error
}

func (s *ApprovalHistoryService) Delete(approvalHistory *model.ApprovalHistory) error {
	return s.db.Delete(approvalHistory).Error
}

func (s *ApprovalHistoryService) GetNotification(userID uint) ([]model.ViewApprovalNotification, error) {
	var approvalNotifications []model.ViewApprovalNotification

	query := s.db.Model(&model.ViewApprovalNotification{})

	if userID != 0 {
		query = query.Where("next_id_user = ?", userID)
	}

	if err := query.Find(&approvalNotifications).Error; err != nil {
		return nil, err
	}

	return approvalNotifications, nil
}
