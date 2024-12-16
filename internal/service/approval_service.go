package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
)

type ApprovalService struct {
	db *gorm.DB
}

func NewApprovalService(db *gorm.DB) *ApprovalService {
	return &ApprovalService{db: db}
}

func (s *ApprovalService) GetByID(approvalID uint) (*model.MstApproval, error) {
	var approval model.MstApproval
	if err := s.db.Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&approval, approvalID).Error; err != nil {
		return nil, err
	}

	return &approval, nil
}

func (s *ApprovalService) GetTotal(search string) (int64, error) {
	var count int64

	query := s.db.Model(&model.MstApproval{})

	if search != "" {
		query = query.Where("status ILIKE ?", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *ApprovalService) GetAll(offset, limit int, search string) ([]model.MstApproval, error) {
	var approvals []model.MstApproval

	query := s.db.Model(&model.MstApproval{}).Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Offset(offset).Limit(limit)

	if search != "" {
		query = query.Where("status ILIKE ?", "%"+search+"%")
	}

	if err := query.Find(&approvals).Error; err != nil {
		return nil, err
	}

	return approvals, nil
}

func (s *ApprovalService) Create(approval []*model.MstApproval) error {
	return s.db.Create(approval).Error
}

func (s *ApprovalService) Update(approval []*model.MstApproval) error {
	return s.db.Save(approval).Error
}

func (s *ApprovalService) Delete(approval *model.MstApproval) error {
	return s.db.Delete(approval).Error
}

func (s *ApprovalService) DeleteByIdMenu(idMenu uint) error {
	return s.db.Where("id_menu = ?", idMenu).Delete(&model.MstApproval{}).Error
}
