package model

import "time"

type MstReason struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	IDMenu      uint       `json:"id_menu,omitempty"`
	Key         string     `json:"key,omitempty"`
	Code        string     `json:"code,omitempty"`
	Description string     `json:"description,omitempty"`
	Remarks     string     `json:"remarks,omitempty"`
	IDCreatedby uint       `json:"id_createdby,omitempty"`
	IDUpdatedby uint       `json:"id_updatedby,omitempty"`
	CreatedAt   *time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt   *time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`

	Menu      *MstMenu `gorm:"foreignKey:ID;references:IDMenu" json:"menu,omitempty"`
	CreatedBy *MstUser `gorm:"foreignKey:ID;references:IDCreatedby" json:"created_by,omitempty"`
	UpdatedBy *MstUser `gorm:"foreignKey:ID;references:IDUpdatedby" json:"updated_by,omitempty"`
}
