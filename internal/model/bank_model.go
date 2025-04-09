package model

import (
	"time"
)

type MstBank struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	AccountNum  string     `json:"account_num"`
	IDAccount   uint       `json:"id_account"`
	IDCurrency  uint       `json:"id_currency"`
	BIC         string     `json:"bic,omitempty"`
	Country     string     `json:"country"`
	State       string     `json:"state"`
	City        string     `json:"city"`
	Address     string     `json:"address"`
	ZipCode     string     `json:"zip_code"`
	Remarks     string     `json:"remarks,omitempty"`
	IDCreatedby uint       `json:"id_createdby,omitempty"`
	IDUpdatedby uint       `json:"id_updatedby,omitempty"`
	CreatedAt   *time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt   *time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`

	Account   *MstChartOfAccount `gorm:"foreignKey:ID;references:IDAccount" json:"account,omitempty"`
	Currency  *MstCurrency       `gorm:"foreignKey:ID;references:IDCurrency" json:"currency,omitempty"`
	CreatedBy *MstUser           `gorm:"foreignKey:ID;references:IDCreatedby" json:"created_by,omitempty"`
	UpdatedBy *MstUser           `gorm:"foreignKey:ID;references:IDUpdatedby" json:"updated_by,omitempty"`
}
