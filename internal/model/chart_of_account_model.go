package model

import (
	"time"
)

type MstChartOfAccount struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	Account          int        `json:"account"`
	Description      string     `json:"description,omitempty"`
	Type             string     `json:"type,omitempty"`
	Class            string     `json:"class,omitempty"`
	ExchangeRateType string     `json:"exchange_rate_type,omitempty"`
	Remarks          string     `json:"remarks,omitempty"`
	IDCreatedby      uint       `json:"id_createdby,omitempty"`
	IDUpdatedby      uint       `json:"id_updatedby,omitempty"`
	CreatedAt        *time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt        *time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`

	CreatedBy *MstUser `gorm:"foreignKey:ID;references:IDCreatedby" json:"created_by,omitempty"`
	UpdatedBy *MstUser `gorm:"foreignKey:ID;references:IDUpdatedby" json:"updated_by,omitempty"`
}
