package model

import (
	"time"
)

type MstMenu struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Label       string     `json:"label,omitempty"`
	Path        string     `json:"path,omitempty"`
	IDParent    uint       `json:"id_parent,omitempty"`
	Icon        string     `json:"icon,omitempty"`
	Sort        uint       `json:"sort,omitempty"`
	IsDelete    uint       `json:"is_delete,omitempty"`
	IDCreatedby uint       `json:"id_createdby,omitempty"`
	IDUpdatedby uint       `json:"id_updatedby,omitempty"`
	CreatedAt   *time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt   *time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`

	CreatedBy     *MstUser       `gorm:"foreignKey:ID;references:IDCreatedby" json:"created_by,omitempty"`
	UpdatedBy     *MstUser       `gorm:"foreignKey:ID;references:IDUpdatedby" json:"updated_by,omitempty"`
	MenuApprovals []*MstApproval `gorm:"foreignKey:IDMenu" json:"menu_approvals,omitempty"`

	Parent          *MstMenu           `gorm:"foreignKey:IDParent" json:"parent,omitempty"`
	Children        []*MstMenu         `gorm:"foreignKey:IDParent" json:"children,omitempty"`
	RolePermissions *MstRolePermission `gorm:"foreignKey:IDMenu;references:ID" json:"role_permissions,omitempty"`
}
