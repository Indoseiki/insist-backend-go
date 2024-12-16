package service

import (
	"insist-backend-golang/internal/model"

	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) GetByID(userID uint) (*model.MstUser, error) {
	var user model.MstUser
	if err := s.db.Preload("Dept", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, code, description")
	}).First(&user, userID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *AuthService) GetByUsername(username string) (*model.MstUser, error) {
	var user model.MstUser
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *AuthService) UpdatePassword(userID uint, newPassword string) error {
	return s.db.Model(&model.MstUser{ID: userID}).Update("password", newPassword).Error
}

func (s *AuthService) UpdateRefreshToken(userID uint, refreshToken string) error {
	return s.db.Model(&model.MstUser{ID: userID}).Update("refresh_token", refreshToken).Error
}

func (s *AuthService) UpdateTwoFactorAuth(userID uint, isTwoFa bool) error {
	return s.db.Model(&model.MstUser{ID: userID}).Update("is_two_fa", isTwoFa).Error
}

func (s *AuthService) UpdateTwoFactorAuthKey(userID uint, otpKey string, otpUrl string) error {
	return s.db.Model(&model.MstUser{ID: userID}).Update("otp_key", otpKey).Update("otp_url", otpUrl).Error
}
