package model

import (
	"time"
)

type MstItemProduct struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	IDItemCategory    uint       `json:"id_item_category"`
	IDItemSubCategory *uint      `json:"id_item_sub_category,omitempty"`
	Code              string     `json:"code"`
	Description       string     `json:"description"`
	Remarks           string     `json:"remarks,omitempty"`
	IDCreatedby       uint       `json:"id_createdby,omitempty"`
	IDUpdatedby       uint       `json:"id_updatedby,omitempty"`
	CreatedAt         *time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt         *time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`

	ItemCategory    *MstItemCategory    `gorm:"foreignKey:ID;references:IDItemCategory" json:"item_category,omitempty"`
	ItemSubCategory *MstItemSubCategory `gorm:"foreignKey:ID;references:IDItemSubCategory" json:"item_sub_category,omitempty"`
	CreatedBy       *MstUser            `gorm:"foreignKey:ID;references:IDCreatedby" json:"created_by,omitempty"`
	UpdatedBy       *MstUser            `gorm:"foreignKey:ID;references:IDUpdatedby" json:"updated_by,omitempty"`
}
