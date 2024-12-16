package model

import (
	"time"
)

type PasswordReset struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	IDUser      int       `json:"id_user"`
	Token       string    `json:"token"`
	IsUsed      bool      `json:"is_used"`
	ExpiredAt   time.Time `json:"expired_at"`
	IDCreatedby uint      `json:"id_createdby"`
	IDUpdatedby uint      `json:"id_updatedby"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
