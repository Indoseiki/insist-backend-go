package model

import "time"

type MstItemRawMaterial struct {
	ID                uint       `gorm:"primaryKey" json:"id"`
	IDItem            uint       `json:"id_item"`
	IDItemProductType uint       `json:"id_item_product_type"`
	IDItemGroupType   uint       `json:"id_item_group_type"`
	IDItemProcess     uint       `json:"id_item_process"`
	IDItemSurface     uint       `json:"id_item_surface"`
	IDItemSource      uint       `json:"id_item_source"`
	IDCreatedby       uint       `json:"id_createdby,omitempty"`
	IDUpdatedby       uint       `json:"id_updatedby,omitempty"`
	CreatedAt         *time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt         *time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`

	Item            *MstItem            `gorm:"foreignKey:ID;references:IDItem" json:"item,omitempty"`
	ItemProductType *MstItemProductType `gorm:"foreignKey:ID;references:IDItemProductType" json:"item_product_type,omitempty"`
	ItemGroupType   *MstItemGroupType   `gorm:"foreignKey:ID;references:IDItemGroupType" json:"item_group_type,omitempty"`
	ItemProcess     *MstItemProcess     `gorm:"foreignKey:ID;references:IDItemProcess" json:"item_process,omitempty"`
	ItemSurface     *MstItemSurface     `gorm:"foreignKey:ID;references:IDItemSurface" json:"item_surface,omitempty"`
	ItemSource      *MstItemSource      `gorm:"foreignKey:ID;references:IDItemSource" json:"item_source,omitempty"`
	CreatedBy       *MstUser            `gorm:"foreignKey:ID;references:IDCreatedby" json:"created_by,omitempty"`
	UpdatedBy       *MstUser            `gorm:"foreignKey:ID;references:IDUpdatedby" json:"updated_by,omitempty"`
}
