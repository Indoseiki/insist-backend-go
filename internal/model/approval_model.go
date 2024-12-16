package model

import (
	"time"
)

type MstApproval struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	IDMenu      uint       `json:"id_menu,omitempty"`
	Status      string     `json:"status,omitempty"`
	Action      string     `json:"action,omitempty"`
	Count       uint       `json:"count,omitempty"`
	Level       uint       `json:"level,omitempty"`
	IDCreatedby uint       `json:"id_createdby,omitempty"`
	IDUpdatedby uint       `json:"id_updatedby,omitempty"`
	CreatedAt   *time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt   *time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`

	CreatedBy     *MstUser           `gorm:"foreignKey:ID;references:IDCreatedby" json:"created_by,omitempty"`
	UpdatedBy     *MstUser           `gorm:"foreignKey:ID;references:IDUpdatedby" json:"updated_by,omitempty"`
	Menu          *MstMenu           `gorm:"foreignKey:IDMenu;references:ID" json:"menu,omitempty"`
	ApprovalUsers []*MstApprovalUser `gorm:"foreignKey:IDApproval" json:"approval_users,omitempty"`
}
