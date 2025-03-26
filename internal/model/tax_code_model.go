package model

import (
	"time"
)

type MstTaxCode struct {
	ID                     uint      `gorm:"primaryKey" json:"id"`
	Name                   string    `json:"name" gorm:"not null"`
	Description            string    `json:"description" gorm:"not null"`
	Type                   string    `json:"type" gorm:"not null"`
	Rate                   float64   `json:"rate" gorm:"not null"`
	IncludePrice           bool      `json:"include_price" gorm:"not null;default:false"`
	IncludeDiscount        bool      `json:"include_discount" gorm:"not null;default:false"`
	IncludeRestockFee      bool      `json:"include_restock_fee" gorm:"not null;default:false"`
	Deductible             bool      `json:"deductible" gorm:"not null;default:false"`
	IncludeFreight         bool      `json:"include_freight" gorm:"not null;default:false"`
	IncludeDuty            bool      `json:"include_duty" gorm:"not null;default:false"`
	IncludeBrokerage       bool      `json:"include_brokerage" gorm:"not null;default:false"`
	IncludeInsurance       bool      `json:"include_insurance" gorm:"not null;default:false"`
	IncludeLocalFreight    bool      `json:"include_local_freight" gorm:"not null;default:false"`
	IncludeMisc            bool      `json:"include_misc" gorm:"not null;default:false"`
	IncludeSurcharge       bool      `json:"include_surcharge" gorm:"not null;default:false"`
	AssessOnReturn         bool      `json:"assess_on_return" gorm:"not null;default:false"`
	IncludeTaxOnPrevSystem bool      `json:"include_tax_on_prev_system" gorm:"not null;default:false"`
	IDAccountAR            uint      `json:"id_account_ar" gorm:"not null"`
	IDAccountARProcess     uint      `json:"id_account_ar_process" gorm:"not null"`
	IDAccountAP            uint      `json:"id_account_ap" gorm:"not null"`
	IDCreatedby            uint      `json:"id_createdby,omitempty"`
	IDUpdatedby            uint      `json:"id_updatedby,omitempty"`
	CreatedAt              time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt              time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	CreatedBy     *MstUser           `gorm:"foreignKey:ID;references:IDCreatedby" json:"created_by,omitempty"`
	UpdatedBy     *MstUser           `gorm:"foreignKey:ID;references:IDUpdatedby" json:"updated_by,omitempty"`
	AccountAR     *MstChartOfAccount `gorm:"foreignKey:ID;references:IDAccountAR" json:"account_ar,omitempty"`
	AccountARProc *MstChartOfAccount `gorm:"foreignKey:ID;references:IDAccountARProcess" json:"account_ar_process,omitempty"`
	AccountAP     *MstChartOfAccount `gorm:"foreignKey:ID;references:IDAccountAP" json:"account_ap,omitempty"`
}
