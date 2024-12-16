package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetByID(userID uint) (*model.MstUser, error) {
	var user model.MstUser
	if err := s.db.Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("Dept", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, code, description")
	}).Omit("password", "otp_key", "otp_url", "refresh_token").First(&user, userID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) GetTotal(search string, idDept uint, isActive string) (int64, error) {
	var count int64

	query := s.db.Model(&model.MstUser{})

	if search != "" {
		query = query.Where("name ILIKE ? OR username ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if idDept != 0 {
		query = query.Where("id_dept = ?", idDept)
	}

	if isActive != "" {
		query = query.Where("is_active = ?", isActive)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (s *UserService) GetAll(offset, limit int, search string, idDept uint, isActive string, sortBy string, sortDirection bool) ([]model.MstUser, error) {
	var users []model.MstUser

	query := s.db.Model(&model.MstUser{}).Preload("CreatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("UpdatedBy", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, name")
	}).Preload("Dept", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, code, description")
	}).Omit("password", "otp_key", "otp_url", "refresh_token").Offset(offset).Limit(limit)

	if sortBy != "" {
		query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: sortBy}, Desc: sortDirection})
	} else {
		query = query.Order("updated_at ASC")
	}

	if search != "" {
		query = query.Where("name ILIKE ? OR username ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if idDept != 0 {
		query = query.Where("id_dept = ?", idDept)
	}

	if isActive != "" {
		query = query.Where("is_active = ?", isActive)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) Create(user *model.MstUser) error {
	return s.db.Create(user).Error
}

func (s *UserService) Update(user *model.MstUser) error {

	var existingUser model.MstUser
	if err := s.db.First(&existingUser, user.ID).Error; err != nil {
		return err
	}

	if !user.IsActive {
		s.db.Model(&model.MstUser{ID: user.ID}).Update("is_active", false)
	}

	if !user.IsTwoFa {
		s.db.Model(&model.MstUser{ID: user.ID}).Update("is_two_fa", false)
	}

	return s.db.Model(&existingUser).Updates(user).Error
}

func (s *UserService) Delete(user *model.MstUser) error {
	return s.db.Delete(user).Error
}

func (s *UserService) UpdatePassword(userID uint, newPassword string) error {
	return s.db.Model(&model.MstUser{ID: userID}).Update("password", newPassword).Error
}
