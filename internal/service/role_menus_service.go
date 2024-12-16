package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
)

type RoleMenuService struct {
	db *gorm.DB
}

func NewRoleMenuService(db *gorm.DB) *RoleMenuService {
	return &RoleMenuService{db: db}
}

func (s *RoleMenuService) GetByID(roleMenuID uint) (*model.MstRole, error) {
	var roleMenu model.MstRole
	if err := s.db.Select("id", "name").Preload("RoleMenus.Menu", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, label")
	}).First(&roleMenu, roleMenuID).Error; err != nil {
		return nil, err
	}

	return &roleMenu, nil
}

func (s *RoleMenuService) GetTotal(search string) (int64, error) {
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

func (s *RoleMenuService) GetAll(offset, limit int, search string) ([]model.MstRole, error) {
	var roleMenus []model.MstRole

	query := s.db.Model(&model.MstRole{}).Preload("RoleMenus.Menu", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, label")
	}).Select("id", "name").Offset(offset).Limit(limit)

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if err := query.Find(&roleMenus).Error; err != nil {
		return nil, err
	}

	return roleMenus, nil
}

func (s *RoleMenuService) Create(roleMenu *[]model.MstRoleMenu) error {
	return s.db.Create(roleMenu).Error
}

func (s *RoleMenuService) Update(roleMenu *model.MstRoleMenu) error {
	return s.db.Save(roleMenu).Error
}

func (s *RoleMenuService) Delete(roleMenuID uint) error {
	return s.db.Where("id_role = ?", roleMenuID).Delete(&model.MstRoleMenu{}).Error
}
