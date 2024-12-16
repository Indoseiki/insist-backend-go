package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
)

type ApprovalUserService struct {
	db *gorm.DB
}

func NewApprovalUserService(db *gorm.DB) *ApprovalUserService {
	return &ApprovalUserService{db: db}
}

func (s *ApprovalUserService) GetTotal(search string) (int64, error) {
	var count int64

	query := s.db.Model(&model.MstApprovalUser{}).Joins("JOIN mst_users u ON u.id = mst_approval_users.id_user")

	if search != "" {
		query.Where("u.name ILIKE ?", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *ApprovalUserService) GetAll(offset, limit int, search string) ([]model.MstApprovalUser, error) {
	var approvalUsers []model.MstApprovalUser

	query := s.db.Model(&model.MstApprovalUser{}).Joins("JOIN mst_users u ON u.id = mst_approval_users.id_user").Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("Approval", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, id_menu, status, action, count, level")
	}).Offset(offset).Limit(limit)

	if search != "" {
		query.Where("u.name ILIKE ?", "%"+search+"%")
	}

	if err := query.Find(&approvalUsers).Error; err != nil {
		return nil, err
	}

	return approvalUsers, nil
}

func (s *ApprovalUserService) Create(approvalUser []*model.MstApprovalUser) error {
	return s.db.Create(approvalUser).Error
}

func (s *ApprovalUserService) Delete(approvalID uint) error {
	return s.db.Where("id_approval = ?", approvalID).Delete(&model.MstApprovalUser{}).Error
}
