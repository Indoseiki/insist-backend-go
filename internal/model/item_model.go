package model

import "time"

type MstItem struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	IDItemCategory   uint       `json:"id_item_category"`
	IDUOM            uint       `json:"id_uom"`
	Code             string     `json:"code"`
	Description      string     `json:"description"`
	InforCode        string     `json:"infor_code"`
	InforDescription string     `json:"infor_description"`
	Remarks          string     `json:"remarks,omitempty"`
	IDCreatedby      uint       `json:"id_createdby,omitempty"`
	IDUpdatedby      uint       `json:"id_updatedby,omitempty"`
	CreatedAt        *time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt        *time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`

	ItemCategory *MstItemCategory `gorm:"foreignKey:ID;references:IDItemCategory" json:"item_category,omitempty"`
	UOM          *MstUoms         `gorm:"foreignKey:ID;references:IDUOM" json:"uom,omitempty"`
	CreatedBy    *MstUser         `gorm:"foreignKey:ID;references:IDCreatedby" json:"created_by,omitempty"`
	UpdatedBy    *MstUser         `gorm:"foreignKey:ID;references:IDUpdatedby" json:"updated_by,omitempty"`
}
