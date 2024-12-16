package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RoleService struct {
	db *gorm.DB
}

func NewRoleService(db *gorm.DB) *RoleService {
	return &RoleService{db: db}
}

func (s *RoleService) GetByID(roleID uint) (*model.MstRole, error) {
	var role model.MstRole
	if err := s.db.Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&role, roleID).Error; err != nil {
		return nil, err
	}

	return &role, nil
}

func (s *RoleService) GetTotal(search string) (int64, error) {
	var count int64

	query := s.db.Model(&model.MstRole{})

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *RoleService) GetAll(offset, limit int, search string, sortBy string, sortDirection bool) ([]model.MstRole, error) {
	var roles []model.MstRole

	query := s.db.Model(&model.MstRole{}).Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
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
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if err := query.Find(&roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}

func (s *RoleService) Create(role *model.MstRole) error {
	return s.db.Create(role).Error
}

func (s *RoleService) Update(role *model.MstRole) error {
	return s.db.Save(role).Error
}

func (s *RoleService) Delete(role *model.MstRole) error {
	return s.db.Delete(role).Error
}
