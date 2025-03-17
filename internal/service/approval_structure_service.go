package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ApprovalStructureService struct {
	db *gorm.DB
}

func NewApprovalStructureService(db *gorm.DB) *ApprovalStructureService {
	return &ApprovalStructureService{db: db}
}

func (s *ApprovalStructureService) GetByID(approvalUserID uint) (*model.MstMenu, error) {
	var approvalUser model.MstMenu
	if err := s.db.Preload("MenuApprovals", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, id_menu, status, action, count, level")
	}).Preload("MenuApprovals.ApprovalUsers.User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&approvalUser, approvalUserID).Error; err != nil {
		return nil, err
	}

	return &approvalUser, nil
}

func (s *ApprovalStructureService) GetTotal(search string) (int64, error) {
	var count int64

	query := s.db.Model(&model.MstMenu{})

	if search != "" {
		query = query.Where("label ILIKE ?", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *ApprovalStructureService) GetAll(offset, limit int, search string, sortBy string, sortDirection bool) ([]model.MstMenu, error) {
	var approvalUsers []model.MstMenu

	query := s.db.Model(&model.MstMenu{}).Preload("MenuApprovals", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, id_menu, status, action, count, level").Order("level ASC")
	}).Preload("MenuApprovals.ApprovalUsers.User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
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
		query = query.Where("label ILIKE ?", "%"+search+"%")
	}

	if err := query.Find(&approvalUsers).Error; err != nil {
		return nil, err
	}

	return approvalUsers, nil
}

func (s *ApprovalStructureService) GetAllByMenu(userID uint, path string) ([]model.ViewApprovalStructure, error) {
	var approvalStructures []model.ViewApprovalStructure

	query := s.db.Model(&model.ViewApprovalStructure{})

	if path != "" {
		query = query.Where("id_user = ? AND path = ?", userID, path)
	}

	if err := query.Find(&approvalStructures).Error; err != nil {
		return nil, err
	}

	return approvalStructures, nil
}
