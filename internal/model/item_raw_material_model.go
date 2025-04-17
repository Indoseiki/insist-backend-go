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
	DiameterSize      string     `json:"diameter_size"`
	LengthSize        string     `json:"length_size"`
	InnerDiameterSize string     `json:"inner_diameter_size"`
	IDCreatedby       uint       `json:"id_createdby,omitempty"`
	IDUpdatedby       uint       `json:"id_updatedby,omitempty"`
	CreatedAt         *time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt         *time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`

	Item            *MstItem            `gorm:"references:ID;foreignKey:IDItem" json:"item,omitempty"`
	ItemProductType *MstItemProductType `gorm:"references:ID;foreignKey:IDItemProductType" json:"item_product_type,omitempty"`
	ItemGroupType   *MstItemGroupType   `gorm:"references:ID;foreignKey:IDItemGroupType" json:"item_group_type,omitempty"`
	ItemProcess     *MstItemProcess     `gorm:"references:ID;foreignKey:IDItemProcess" json:"item_process,omitempty"`
	ItemSurface     *MstItemSurface     `gorm:"references:ID;foreignKey:IDItemSurface" json:"item_surface,omitempty"`
	ItemSource      *MstItemSource      `gorm:"references:ID;foreignKey:IDItemSource" json:"item_source,omitempty"`
	CreatedBy       *MstUser            `gorm:"foreignKey:ID;references:IDCreatedby" json:"created_by,omitempty"`
	UpdatedBy       *MstUser            `gorm:"foreignKey:ID;references:IDUpdatedby" json:"updated_by,omitempty"`
}
