package model

import (
	"time"
)

type MstLocation struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	IDWarehouse uint       `json:"id_warehouse,omitempty"`
	Location    string     `json:"location,omitempty"`
	Remarks     string     `json:"remarks,omitempty"`
	IDCreatedby uint       `json:"id_createdby,omitempty"`
	IDUpdatedby uint       `json:"id_updatedby,omitempty"`
	CreatedAt   *time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt   *time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`

	CreatedBy *MstUser `gorm:"foreignKey:ID;references:IDCreatedby" json:"created_by,omitempty"`
	UpdatedBy *MstUser `gorm:"foreignKey:ID;references:IDUpdatedby" json:"updated_by,omitempty"`

	Warehouse *MstBuilding `gorm:"foreignKey:ID;references:IDWarehouse" json:"warehouse,omitempty"`
}
