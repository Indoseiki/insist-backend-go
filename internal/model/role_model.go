package model

import (
	"time"
)

type MstRole struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Name        string     `json:"name,omitempty"`
	IDCreatedby uint       `json:"id_createdby,omitempty"`
	IDUpdatedby uint       `json:"id_updatedby,omitempty"`
	CreatedAt   *time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt   *time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`

	CreatedBy *MstUser `gorm:"foreignKey:ID;references:IDCreatedby" json:"created_by,omitempty"`
	UpdatedBy *MstUser `gorm:"foreignKey:ID;references:IDUpdatedby" json:"updated_by,omitempty"`

	RoleMenus       []*MstRoleMenu       `gorm:"foreignKey:IDRole;references:ID" json:"role_menus,omitempty"`
	RolePermissions []*MstRolePermission `gorm:"foreignKey:IDRole;references:ID" json:"role_permissions,omitempty"`
}
