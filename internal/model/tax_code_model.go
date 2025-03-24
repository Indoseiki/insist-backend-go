package model

import (
	"time"
)

type MstTaxCode struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Name            string    `json:"name" gorm:"not null"`
	Description     string    `json:"description" gorm:"not null"`
	Type            string    `json:"type" gorm:"not null"`
	Rate            float64   `json:"rate" gorm:"not null"`
	Include         string    `json:"include" gorm:"not null"`
	IDAccountAR     uint      `json:"id_account_ar" gorm:"not null"`
	IDAccountARProc uint      `json:"id_account_ar_process" gorm:"not null"`
	IDAccountAP     uint      `json:"id_account_ap" gorm:"not null"`
	IDCreatedby     uint      `json:"id_createdby"`
	IDUpdatedby     uint      `json:"id_updatedby"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	AccountAR     *MstChartOfAccount `gorm:"foreignKey:ID;references:IDAccountAR" json:"account_ar,omitempty"`
	AccountARProc *MstChartOfAccount `gorm:"foreignKey:ID;references:IDAccountARProc" json:"account_ar_process,omitempty"`
	AccountAP     *MstChartOfAccount `gorm:"foreignKey:ID;references:IDAccountAP" json:"account_ap,omitempty"`
	CreatedBy     *MstUser           `gorm:"foreignKey:ID;references:IDCreatedby" json:"created_by,omitempty"`
	UpdatedBy     *MstUser           `gorm:"foreignKey:ID;references:IDUpdatedby" json:"updated_by,omitempty"`
}
