package model

import (
	"time"
)

type MstCurrencyRate struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	IDFromCurrency uint       `json:"id_from_currency,omitempty"`
	IDToCurrency   uint       `json:"id_to_currency,omitempty"`
	BuyRate        float64    `json:"buy_rate,omitempty"`
	SellRate       float64    `json:"sell_rate,omitempty"`
	EffectiveDate  time.Time  `json:"effective_date,omitempty"`
	IDCreatedby    uint       `json:"id_createdby,omitempty"`
	IDUpdatedby    uint       `json:"id_updatedby,omitempty"`
	CreatedAt      *time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt      *time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`

	FromCurrency *MstCurrency `gorm:"foreignKey:ID;references:IDFromCurrency" json:"from_currency,omitempty"`
	ToCurrency   *MstCurrency `gorm:"foreignKey:ID;references:IDToCurrency" json:"to_currency,omitempty"`
	CreatedBy    *MstUser     `gorm:"foreignKey:ID;references:IDCreatedby" json:"created_by,omitempty"`
	UpdatedBy    *MstUser     `gorm:"foreignKey:ID;references:IDUpdatedby" json:"updated_by,omitempty"`
}
