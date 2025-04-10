package model

import (
	"time"
)

type MstBillingTerm struct {
	ID                        uint       `gorm:"primaryKey" json:"id"`
	Code                      string     `json:"code"`
	Description               string     `json:"description"`
	DueDays                   int        `json:"due_days,omitempty"`
	DiscountDays              int        `json:"discount_days,omitempty"`
	IsCashOnly                bool       `json:"is_cash_only,omitempty"`
	ProxDueDay                int        `json:"prox_due_day,omitempty"`
	ProxDiscountDay           int        `json:"prox_discount_day,omitempty"`
	ProxMonthsForward         int        `json:"prox_months_forward,omitempty"`
	ProxDiscountMonthsForward int        `json:"prox_discount_months_forward,omitempty"`
	CutoffDay                 int        `json:"cutoff_day,omitempty"`
	DiscountPercent           float64    `gorm:"type:decimal(5,3)" json:"discount_percent,omitempty"`
	HolidayOffsetMethod       string     `json:"holiday_offset_method,omitempty"`
	IsAdvancedTerms           bool       `json:"is_advanced_terms,omitempty"`
	ProxCode                  int        `json:"prox_code,omitempty"`
	IDCreatedby               uint       `json:"id_createdby,omitempty"`
	IDUpdatedby               uint       `json:"id_updatedby,omitempty"`
	CreatedAt                 *time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt                 *time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`

	CreatedBy *MstUser `gorm:"foreignKey:ID;references:IDCreatedby" json:"created_by,omitempty"`
	UpdatedBy *MstUser `gorm:"foreignKey:ID;references:IDUpdatedby" json:"updated_by,omitempty"`
}
