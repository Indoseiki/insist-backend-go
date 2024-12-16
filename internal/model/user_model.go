package model

import (
	"time"
)

type MstUser struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	IDDept       uint       `json:"id_dept,omitempty"`
	Name         string     `json:"name,omitempty"`
	Email        string     `json:"email,omitempty"`
	Username     string     `json:"username,omitempty"`
	Password     string     `json:"password,omitempty"`
	RefreshToken string     `json:"refresh_token,omitempty"`
	OtpKey       string     `json:"otp_key,omitempty"`
	OtpUrl       string     `json:"otp_url,omitempty"`
	IsActive     bool       `json:"is_active"`
	IsTwoFa      bool       `json:"is_two_fa"`
	IDCreatedby  uint       `json:"id_createdby,omitempty"`
	IDUpdatedby  uint       `json:"id_updatedby,omitempty"`
	CreatedAt    *time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt    *time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`

	Dept      *MstDept       `gorm:"foreignKey:IDDept;references:ID" json:"dept,omitempty"`
	CreatedBy *MstUser       `gorm:"foreignKey:IDCreatedby;references:ID" json:"created_by,omitempty"`
	UpdatedBy *MstUser       `gorm:"foreignKey:IDUpdatedby;references:ID" json:"updated_by,omitempty"`
	UserRoles []*MstUserRole `gorm:"foreignKey:IDUser;references:ID" json:"user_roles,omitempty"`
}
