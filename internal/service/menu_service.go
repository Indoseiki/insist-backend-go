package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MenuService struct {
	db *gorm.DB
}

func NewMenuService(db *gorm.DB) *MenuService {
	return &MenuService{db: db}
}

func (s *MenuService) GetByID(menuID uint) (*model.MstMenu, error) {
	var menu model.MstMenu
	if err := s.db.Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&menu, menuID).Error; err != nil {
		return nil, err
	}

	return &menu, nil
}

func (s *MenuService) GetTotalByUser(userID uint, search string) (int64, error) {
	var count int64

	query := s.db.Table("mst_menus").
		Select("COUNT(DISTINCT mst_menus.id)").
		Joins("JOIN mst_role_menus ON mst_role_menus.id_menu = mst_menus.id").
		Joins("JOIN mst_user_roles ON mst_user_roles.id_role = mst_role_menus.id_role").
		Where("mst_user_roles.id_user = ? AND mst_menus.is_delete = ?", userID, 0)

	if search != "" {
		query = query.Where("mst_menus.label ILIKE ?", "%"+search+"%")
	}

	// Hitung total
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *MenuService) GetByUser(userID uint, offset, limit int, search string) ([]model.MstMenu, error) {
	var menus []model.MstMenu

	query := s.db.Table("mst_menus").
		Select("mst_menus.id, mst_menus.label, mst_menus.path, mst_menus.id_parent, mst_menus.icon").
		Joins("JOIN mst_role_menus ON mst_role_menus.id_menu = mst_menus.id").
		Joins("JOIN mst_user_roles ON mst_user_roles.id_role = mst_role_menus.id_role").
		Where("mst_user_roles.id_user = ? AND mst_menus.is_delete = ?", userID, 0).
		Group("mst_menus.id, mst_menus.label, mst_menus.path, mst_menus.id_parent, mst_menus.icon").
		Order("mst_menus.label ASC").
		Offset(offset).
		Limit(limit)

	if search != "" {
		query = query.Where("mst_menus.label ILIKE ?", "%"+search+"%")
	}

	if err := query.Find(&menus).Error; err != nil {
		return nil, err
	}

	return menus, nil
}

func (s *MenuService) GetTotal(search string) (int64, error) {
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

func (s *MenuService) GetAll(offset, limit int, search string, sortBy string, sortDirection bool) ([]model.MstMenu, error) {
	var menus []model.MstMenu

	query := s.db.Model(&model.MstMenu{}).Preload("Parent", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, label")
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

	if err := query.Find(&menus).Error; err != nil {
		return nil, err
	}

	return menus, nil
}

func (s *MenuService) Create(menu *model.MstMenu) error {
	return s.db.Create(menu).Error
}

func (s *MenuService) Update(menu *model.MstMenu) error {
	return s.db.Save(menu).Error
}

func (s *MenuService) Delete(menu *model.MstMenu) error {
	return s.db.Delete(menu).Error
}

func (s *MenuService) GetMenuTree() ([]model.MstMenu, error) {
	var rootMenus []model.MstMenu

	if err := s.db.Where("id_parent = ? AND is_delete = ?", 0, 0).Order("sort ASC").Order("label ASC").Find(&rootMenus).Error; err != nil {
		return nil, err
	}

	for i := range rootMenus {
		if err := s.loadChildren(&rootMenus[i]); err != nil {
			return nil, err
		}
	}

	return rootMenus, nil
}

func (s *MenuService) loadChildren(menu *model.MstMenu) error {
	var children []model.MstMenu

	if err := s.db.Where("id_parent = ? AND is_delete = ?", menu.ID, 0).Order("sort ASC").Order("label ASC").Find(&children).Error; err != nil {
		return err
	}

	var childrenPointers []*model.MstMenu
	for i := range children {
		childrenPointers = append(childrenPointers, &children[i])
	}

	menu.Children = childrenPointers

	for _, child := range menu.Children {
		if err := s.loadChildren(child); err != nil {
			return err
		}
	}

	return nil
}

func (s *MenuService) GetMenuTreeByUser(userID uint) ([]model.MstMenu, error) {
	var rootMenus []model.MstMenu

	var roleMenus []model.MstRoleMenu
	if err := s.db.Table("mst_role_menus").
		Joins("JOIN mst_user_roles ON mst_user_roles.id_role = mst_role_menus.id_role").
		Where("mst_user_roles.id_user = ?", userID).
		Find(&roleMenus).Error; err != nil {
		return nil, err
	}

	var roleMenuIDs []uint
	for _, roleMenu := range roleMenus {
		roleMenuIDs = append(roleMenuIDs, roleMenu.IDMenu)
	}

	if err := s.db.Where("id_parent = ? AND is_delete = ? AND id IN (?)", 0, 0, roleMenuIDs).Order("sort ASC").Order("label ASC").Find(&rootMenus).Error; err != nil {
		return nil, err
	}

	for i := range rootMenus {
		if err := s.loadChildrenByUser(&rootMenus[i], roleMenuIDs); err != nil {
			return nil, err
		}
	}

	return rootMenus, nil
}

func (s *MenuService) loadChildrenByUser(menu *model.MstMenu, roleMenuIDs []uint) error {
	var children []model.MstMenu

	if err := s.db.Where("id_parent = ? AND is_delete = ? AND id IN (?)", menu.ID, 0, roleMenuIDs).Order("sort ASC").Order("label ASC").Find(&children).Error; err != nil {
		return err
	}

	var childrenPointers []*model.MstMenu
	for i := range children {
		childrenPointers = append(childrenPointers, &children[i])
	}

	menu.Children = childrenPointers

	for _, child := range menu.Children {
		if err := s.loadChildrenByUser(child, roleMenuIDs); err != nil {
			return err
		}
	}

	return nil
}
