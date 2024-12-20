package service

import (
	"insist-backend-golang/internal/dto"
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
)

type RolePermissionService struct {
	db *gorm.DB
}

func NewRolePermissionService(db *gorm.DB) *RolePermissionService {
	return &RolePermissionService{db: db}
}

func (s *RolePermissionService) GetByPath(userID uint, path string) ([]dto.MenuWithPermissions, error) {
	var menus []dto.MenuWithPermissions

	query := s.db.Table("mst_menus").
		Select("mst_menus.label, mst_menus.path, mst_role_permissions.*").
		Joins("JOIN mst_role_permissions ON mst_role_permissions.id_menu = mst_menus.id").
		Joins("JOIN mst_user_roles ON mst_user_roles.id_role = mst_role_permissions.id_role").
		Where("mst_user_roles.id_user = ? AND mst_menus.path = ?", userID, path)

	if err := query.Scan(&menus).Error; err != nil {
		return nil, err
	}

	return menus, nil
}

func (s *RolePermissionService) GetMenuTreeByRole(roleID uint) ([]model.MstMenu, error) {
	var rootMenus []model.MstMenu

	var roleMenus []model.MstRoleMenu
	if err := s.db.Table("mst_role_menus").
		Where("id_role = ?", roleID).
		Find(&roleMenus).Error; err != nil {
		return nil, err
	}

	var roleMenuIDs []uint
	for _, roleMenu := range roleMenus {
		roleMenuIDs = append(roleMenuIDs, roleMenu.IDMenu)
	}

	if err := s.db.Where("id_parent = ? AND id IN (?)", 0, roleMenuIDs).
		Order("sort ASC").
		Preload("RolePermissions", "id_role = ?", roleID).
		Find(&rootMenus).Error; err != nil {
		return nil, err
	}

	for i := range rootMenus {
		if err := s.loadChildrenByRole(&rootMenus[i], roleMenuIDs, roleID); err != nil {
			return nil, err
		}
	}

	return rootMenus, nil
}

func (s *RolePermissionService) loadChildrenByRole(menu *model.MstMenu, roleMenuIDs []uint, roleID uint) error {
	var children []model.MstMenu

	if err := s.db.Where("id_parent = ? AND id IN (?)", menu.ID, roleMenuIDs).
		Order("sort ASC").
		Preload("RolePermissions", "id_role = ?", roleID).
		Find(&children).Error; err != nil {
		return err
	}

	var childrenPointers []*model.MstMenu
	for i := range children {
		childrenPointers = append(childrenPointers, &children[i])
	}

	menu.Children = childrenPointers

	for _, child := range menu.Children {
		if err := s.loadChildrenByRole(child, roleMenuIDs, roleID); err != nil {
			return err
		}
	}

	return nil
}

func (s *RolePermissionService) UpdateOrCreateRolePermission(rolePermission *model.MstRolePermission) error {
	var existingPermission model.MstRolePermission
	err := s.db.Where("id_role = ? AND id_menu = ?", rolePermission.IDRole, rolePermission.IDMenu).First(&existingPermission).Error

	if err == nil {
		if !rolePermission.IsCreate && !rolePermission.IsUpdate && !rolePermission.IsDelete {
			return s.DeleteRolePermission(rolePermission.IDRole, rolePermission.IDMenu)
		}

		existingPermission.IsCreate = rolePermission.IsCreate
		existingPermission.IsUpdate = rolePermission.IsUpdate
		existingPermission.IsDelete = rolePermission.IsDelete

		if err := s.db.Model(&existingPermission).
			Where("id_role = ? AND id_menu = ?", rolePermission.IDRole, rolePermission.IDMenu).
			Updates(map[string]interface{}{
				"is_create": rolePermission.IsCreate,
				"is_update": rolePermission.IsUpdate,
				"is_delete": rolePermission.IsDelete,
			}).Error; err != nil {
			return err
		}
		return nil
	}

	if err == gorm.ErrRecordNotFound {
		if !rolePermission.IsCreate && !rolePermission.IsUpdate && !rolePermission.IsDelete {
			return nil
		}

		if err := s.db.Create(rolePermission).Error; err != nil {
			return err
		}
		return nil
	}

	return err
}

func (s *RolePermissionService) DeleteRolePermission(idRole uint, idMenu uint) error {
	var rolePermission model.MstRolePermission
	err := s.db.Where("id_role = ? AND id_menu = ?", idRole, idMenu).First(&rolePermission).Error
	if err != nil {
		return err
	}

	if err := s.db.Where("id_role = ? AND id_menu = ?", idRole, idMenu).Delete(&rolePermission).Error; err != nil {
		return err
	}

	return nil
}
