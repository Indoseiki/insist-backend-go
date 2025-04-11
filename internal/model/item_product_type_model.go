package model

import (
	"time"
)

type MstItemProductType struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	IDItemProduct uint       `json:"id_item_product"`
	Code          string     `json:"code"`
	Description   string     `json:"description"`
	Remarks       string     `json:"remarks,omitempty"`
	IDCreatedby   uint       `json:"id_createdby,omitempty"`
	IDUpdatedby   uint       `json:"id_updatedby,omitempty"`
	CreatedAt     *time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt     *time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`

	ItemProduct *MstItemProduct `gorm:"foreignKey:ID;references:IDItemProduct" json:"item_product,omitempty"`
	CreatedBy   *MstUser        `gorm:"foreignKey:ID;references:IDCreatedby" json:"created_by,omitempty"`
	UpdatedBy   *MstUser        `gorm:"foreignKey:ID;references:IDUpdatedby" json:"updated_by,omitempty"`
}
