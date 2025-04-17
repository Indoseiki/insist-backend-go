package model

import "time"

type MstMaterialDetail struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	IDMaterial      uint       `json:"id_material"`
	RevNo           int        `json:"rev_no"`
	RmssNum         string     `json:"rmss_num"`
	OdTolerancePlus float64    `json:"od_tolerance_plus"`
	OdToleranceMin  float64    `json:"od_tolerance_min"`
	IdTolerancePlus float64    `json:"id_tolerance_plus"`
	IdToleranceMin  float64    `json:"id_tolerance_min"`
	Width           float64    `json:"width"`
	Height          float64    `json:"height"`
	Ovality         string     `json:"ovality,omitempty"`
	CuttingLength   string     `json:"cutting_length,omitempty"`
	Hardness        string     `json:"hardness,omitempty"`
	CompotitionC    string     `json:"compotition_c,omitempty"`
	CompotitionSi   string     `json:"compotition_si,omitempty"`
	CompotitionMn   string     `json:"compotition_mn,omitempty"`
	CompotitionP    string     `json:"compotition_p,omitempty"`
	CompotitionS    string     `json:"compotition_s,omitempty"`
	CompotitionCu   string     `json:"compotition_cu,omitempty"`
	CompotitionNi   string     `json:"compotition_ni,omitempty"`
	CompotitionCr   string     `json:"compotition_cr,omitempty"`
	CompotitionMo   string     `json:"compotition_mo,omitempty"`
	TensileStrength string     `json:"tensile_strength,omitempty"`
	SaRatio         string     `json:"sa_ratio,omitempty"`
	Origin          string     `json:"origin,omitempty"`
	Remarks         string     `json:"remarks,omitempty"`
	IDCreatedby     uint       `json:"id_createdby,omitempty"`
	IDUpdatedby     uint       `json:"id_updatedby,omitempty"`
	CreatedAt       *time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt       *time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`

	Material  *MstMaterial `gorm:"foreignKey:IDMaterial;references:ID" json:"material,omitempty"`
	CreatedBy *MstUser     `gorm:"foreignKey:ID;references:IDCreatedby" json:"created_by,omitempty"`
	UpdatedBy *MstUser     `gorm:"foreignKey:ID;references:IDUpdatedby" json:"updated_by,omitempty"`
}
