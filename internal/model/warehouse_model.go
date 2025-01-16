package model

import (
	"time"
)

type MstWarehouse struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	IDBuilding  uint       `json:"id_building,omitempty"`
	Location    string     `json:"location,omitempty"`
	Remarks     string     `json:"remarks,omitempty"`
	IDCreatedby uint       `json:"id_createdby,omitempty"`
	IDUpdatedby uint       `json:"id_updatedby,omitempty"`
	CreatedAt   *time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt   *time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`

	CreatedBy *MstUser `gorm:"foreignKey:ID;references:IDCreatedby" json:"created_by,omitempty"`
	UpdatedBy *MstUser `gorm:"foreignKey:ID;references:IDUpdatedby" json:"updated_by,omitempty"`

	Building *MstBuilding `gorm:"foreignKey:ID;references:IDBuilding" json:"building,omitempty"`
}
