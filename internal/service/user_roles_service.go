package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
)

type UserRoleService struct {
	db *gorm.DB
}

func NewUserRoleService(db *gorm.DB) *UserRoleService {
	return &UserRoleService{db: db}
}

func (s *UserRoleService) GetByID(userID uint) (*model.MstUser, error) {
	var user model.MstUser
	if err := s.db.Select("id", "name").Preload("UserRoles.Role", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).First(&user, userID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserRoleService) GetTotal(search string) (int64, error) {
	var count int64

	query := s.db.Model(&model.MstUser{})

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *UserRoleService) GetAll(offset, limit int, search string) ([]model.MstUser, error) {
	var users []model.MstUser

	query := s.db.Model(&model.MstUser{}).Preload("UserRoles.Role", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Select("id", "name").Offset(offset).Limit(limit)

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserRoleService) Create(user *[]model.MstUserRole) error {
	return s.db.Create(user).Error
}

func (s *UserRoleService) Update(user *model.MstUserRole) error {
	return s.db.Save(user).Error
}

func (s *UserRoleService) Delete(userID uint) error {
	return s.db.Where("id_user = ?", userID).Delete(&model.MstUserRole{}).Error
}
