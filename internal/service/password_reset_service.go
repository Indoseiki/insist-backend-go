package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
)

type PasswordResetService struct {
	db *gorm.DB
}

func NewPasswordResetService(db *gorm.DB) *PasswordResetService {
	return &PasswordResetService{db: db}
}

func (s *PasswordResetService) GetPasswordResetByToken(token string) (*model.PasswordReset, error) {
	var passwordReset model.PasswordReset
	if err := s.db.Where("token = ?", token).First(&passwordReset).Error; err != nil {
		return nil, err
	}

	return &passwordReset, nil
}

func (s *PasswordResetService) CreatePasswordReset(passwordReset *model.PasswordReset) error {
	return s.db.Create(passwordReset).Error
}

func (s *PasswordResetService) UpdateUsed(userID uint, isUsed bool) error {
	return s.db.Model(&model.PasswordReset{ID: userID}).Update("is_used", isUsed).Error
}
