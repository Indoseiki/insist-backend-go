package model

import (
	"time"
)

type MstMachine struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	IDCreatedby uint      `json:"id_createdby,omitempty"`
	IDUpdatedby uint      `json:"id_updatedby,omitempty"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`

	CreatedBy *MstUser `gorm:"foreignKey:ID;references:IDCreatedby" json:"created_by,omitempty"`
	UpdatedBy *MstUser `gorm:"foreignKey:ID;references:IDUpdatedby" json:"updated_by,omitempty"`
}

type MstMachineDetail struct {
	ID                  uint       `gorm:"primaryKey" json:"id"`
	IDMachine           uint       `json:"id_machine"`
	RevNo               int        `json:"rev_no"`
	Code                string     `json:"code"`
	CodeOld             *string    `json:"code_old,omitempty"`
	AssetNum            string     `json:"asset_num"`
	AssetNumOld         *string    `json:"asset_num_old,omitempty"`
	Description         string     `json:"description"`
	Name                string     `json:"name"`
	Maker               string     `json:"maker"`
	Power               float64    `json:"power"`
	IDPowerUOM          uint       `json:"id_power_uom"`
	Electricity         float64    `json:"electricity"`
	IDElectricityUOM    uint       `json:"id_electricity_uom"`
	Cavity              int        `json:"cavity"`
	Lubricant           string     `json:"lubricant"`
	LubricantCapacity   float64    `json:"lubricant_capacity"`
	IDLubricantUOM      uint       `json:"id_lubricant_uom"`
	Sliding             string     `json:"sliding"`
	SlidingCapacity     float64    `json:"sliding_capacity"`
	IDSlidingUOM        uint       `json:"id_sliding_uom"`
	Coolant             string     `json:"coolant"`
	CoolantCapacity     float64    `json:"coolant_capacity"`
	IDCoolantUOM        uint       `json:"id_coolant_uom"`
	Hydraulic           string     `json:"hydraulic"`
	HydraulicCapacity   float64    `json:"hydraulic_capacity"`
	IDHydraulicUOM      uint       `json:"id_hydraulic_uom"`
	DimensionFront      float64    `json:"dimension_front"`
	IDDimensionFrontUOM uint       `json:"id_dimension_front_uom"`
	DimensionSide       float64    `json:"dimension_side"`
	IDDimensionSideUOM  uint       `json:"id_dimension_side_uom"`
	IDCreatedby         uint       `json:"id_createdby"`
	IDUpdatedby         uint       `json:"id_updatedby"`
	CreatedAt           *time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt           *time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`

	Machine           *MstMachine `gorm:"foreignKey:ID;references:IDMachine" json:"machine,omitempty"`
	PowerUOM          *MstUoms    `gorm:"foreignKey:ID;references:IDPowerUOM" json:"power_uom,omitempty"`
	ElectricityUOM    *MstUoms    `gorm:"foreignKey:ID;references:IDElectricityUOM" json:"electricity_uom,omitempty"`
	LubricantUOM      *MstUoms    `gorm:"foreignKey:ID;references:IDLubricantUOM" json:"lubricant_uom,omitempty"`
	SlidingUOM        *MstUoms    `gorm:"foreignKey:ID;references:IDLubricantUOM" json:"sliding_uom,omitempty"`
	CoolantUOM        *MstUoms    `gorm:"foreignKey:ID;references:IDCoolantUOM" json:"coolant_uom,omitempty"`
	HydraulicUOM      *MstUoms    `gorm:"foreignKey:ID;references:IDHydraulicUOM" json:"hydraulic_uom,omitempty"`
	DimensionFrontUOM *MstUoms    `gorm:"foreignKey:ID;references:IDDimensionFrontUOM" json:"dimension_front_uom,omitempty"`
	DimensionSideUOM  *MstUoms    `gorm:"foreignKey:ID;references:IDDimensionSideUOM" json:"dimension_side_uom,omitempty"`
	CreatedBy         *MstUser    `gorm:"foreignKey:ID;references:IDCreatedby" json:"created_by,omitempty"`
	UpdatedBy         *MstUser    `gorm:"foreignKey:ID;references:IDUpdatedby" json:"updated_by,omitempty"`
}

type MstMachineStatus struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	IDMachine   uint       `json:"id_machine"`
	IDReason    uint       `json:"id_reason"`
	Remarks     *string    `json:"remarks,omitempty"`
	IDCreatedby uint       `json:"id_createdby"`
	IDUpdatedby uint       `json:"id_updatedby"`
	CreatedAt   *time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt   *time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`

	Machine   *MstMachine `gorm:"foreignKey:ID;references:IDMachine" json:"machine,omitempty"`
	Reason    *MstReason  `gorm:"foreignKey:ID;references:IDReason" json:"reason,omitempty"`
	CreatedBy *MstUser    `gorm:"foreignKey:ID;references:IDCreatedby" json:"created_by,omitempty"`
	UpdatedBy *MstUser    `gorm:"foreignKey:ID;references:IDUpdatedby" json:"updated_by,omitempty"`
}

type ViewMstMachine struct {
	ID                   uint      `json:"id"`
	IDCreatedby          uint      `json:"id_createdby"`
	MachineCreatedbyName string    `json:"machine_createdby_name"`
	IDUpdatedby          uint      `json:"id_updatedby"`
	MachineUpdatedbyName string    `json:"machine_updatedby_name"`
	MachineCreatedAt     time.Time `json:"machine_created_at"`
	MachineUpdatedAt     time.Time `json:"machine_updated_at"`

	DetailID            uint    `json:"detail_id"`
	RevNo               int     `json:"rev_no"`
	Code                string  `json:"code"`
	CodeOld             string  `json:"code_old"`
	AssetNum            string  `json:"asset_num"`
	AssetNumOld         string  `json:"asset_num_old"`
	Description         string  `json:"description"`
	Name                string  `json:"name"`
	Maker               string  `json:"maker"`
	Power               float64 `json:"power"`
	IDPowerUOM          uint    `json:"id_power_uom"`
	PowerUOMCode        string  `json:"power_uom_code"`
	PowerUOMDescription string  `json:"power_uom_description"`

	Electricity               float64 `json:"electricity"`
	IDElectricityUOM          uint    `json:"id_electricity_uom"`
	ElectricityUOMCode        string  `json:"electricity_uom_code"`
	ElectricityUOMDescription string  `json:"electricity_uom_description"`

	Cavity                  int     `json:"cavity"`
	Lubricant               string  `json:"lubricant"`
	LubricantCapacity       float64 `json:"lubricant_capacity"`
	IDLubricantUOM          uint    `json:"id_lubricant_uom"`
	LubricantUOMCode        string  `json:"lubricant_uom_code"`
	LubricantUOMDescription string  `json:"lubricant_uom_description"`

	Sliding               string  `json:"sliding"`
	SlidingCapacity       float64 `json:"sliding_capacity"`
	IDSlidingUOM          uint    `json:"id_sliding_uom"`
	SlidingUOMCode        string  `json:"sliding_uom_code"`
	SlidingUOMDescription string  `json:"sliding_uom_description"`

	Coolant               string  `json:"coolant"`
	CoolantCapacity       float64 `json:"coolant_capacity"`
	IDCoolantUOM          uint    `json:"id_coolant_uom"`
	CoolantUOMCode        string  `json:"coolant_uom_code"`
	CoolantUOMDescription string  `json:"coolant_uom_description"`

	Hydraulic               string  `json:"hydraulic"`
	HydraulicCapacity       float64 `json:"hydraulic_capacity"`
	IDHydraulicUOM          uint    `json:"id_hydraulic_uom"`
	HydraulicUOMCode        string  `json:"hydraulic_uom_code"`
	HydraulicUOMDescription string  `json:"hydraulic_uom_description"`

	DimensionFront               float64 `json:"dimension_front"`
	IDDimensionFrontUOM          uint    `json:"id_dimension_front_uom"`
	DimensionFrontUOMCode        string  `json:"dimension_front_uom_code"`
	DimensionFrontUOMDescription string  `json:"dimension_front_uom_description"`

	DimensionSide               float64 `json:"dimension_side"`
	IDDimensionSideUOM          uint    `json:"id_dimension_side_uom"`
	DimensionSideUOMCode        string  `json:"dimension_side_uom_code"`
	DimensionSideUOMDescription string  `json:"dimension_side_uom_description"`

	DetailIDCreatedby   uint      `json:"detail_id_createdby"`
	DetailCreatedbyName string    `json:"detail_createdby_name"`
	DetailIDUpdatedby   uint      `json:"detail_id_updatedby"`
	DetailUpdatedbyName string    `json:"detail_updatedby_name"`
	DetailCreatedAt     time.Time `json:"detail_created_at"`
	DetailUpdatedAt     time.Time `json:"detail_updated_at"`

	IDReason        uint      `json:"id_reason"`
	Remarks         string    `json:"remarks"`
	StatusCreatedAt time.Time `json:"status_created_at"`
	StatusUpdatedAt time.Time `json:"status_updated_at"`

	ReasonKey         string `json:"reason_key"`
	ReasonCode        string `json:"reason_code"`
	ReasonDescription string `json:"reason_description"`
	ReasonRemarks     string `json:"reason_remarks"`

	ApprovalID            uint      `json:"approval_id"`
	RefTable              string    `json:"ref_table"`
	RefID                 uint      `json:"ref_id"`
	ApprovalKey           string    `json:"approval_key"`
	ApprovalMessage       string    `json:"approval_message"`
	ApprovalCreatedby     uint      `json:"approval_createdby"`
	ApprovalCreatedbyName string    `json:"approval_createdby_name"`
	ApprovalCreatedAt     time.Time `json:"approval_created_at"`
	ApprovalStatus        string    `json:"approval_status"`
	ApprovalAction        string    `json:"approval_action"`
	ApprovalCount         int       `json:"approval_count"`
	ApprovalLevel         int       `json:"approval_level"`
}

type ViewMstMachineDetail struct {
	ID                  uint    `json:"id"`
	IDMachine           uint    `json:"id_machine"`
	RevNo               int     `json:"rev_no"`
	Code                string  `json:"code"`
	CodeOld             string  `json:"code_old"`
	AssetNum            string  `json:"asset_num"`
	AssetNumOld         string  `json:"asset_num_old"`
	Description         string  `json:"description"`
	Name                string  `json:"name"`
	Maker               string  `json:"maker"`
	Power               float64 `json:"power"`
	IDPowerUOM          uint    `json:"id_power_uom"`
	PowerUOMCode        string  `json:"power_uom_code"`
	PowerUOMDescription string  `json:"power_uom_description"`

	Electricity               float64 `json:"electricity"`
	IDElectricityUOM          uint    `json:"id_electricity_uom"`
	ElectricityUOMCode        string  `json:"electricity_uom_code"`
	ElectricityUOMDescription string  `json:"electricity_uom_description"`

	Cavity                  int     `json:"cavity"`
	Lubricant               string  `json:"lubricant"`
	LubricantCapacity       float64 `json:"lubricant_capacity"`
	IDLubricantUOM          uint    `json:"id_lubricant_uom"`
	LubricantUOMCode        string  `json:"lubricant_uom_code"`
	LubricantUOMDescription string  `json:"lubricant_uom_description"`

	Sliding               string  `json:"sliding"`
	SlidingCapacity       float64 `json:"sliding_capacity"`
	IDSlidingUOM          uint    `json:"id_sliding_uom"`
	SlidingUOMCode        string  `json:"sliding_uom_code"`
	SlidingUOMDescription string  `json:"sliding_uom_description"`

	Coolant               string  `json:"coolant"`
	CoolantCapacity       float64 `json:"coolant_capacity"`
	IDCoolantUOM          uint    `json:"id_coolant_uom"`
	CoolantUOMCode        string  `json:"coolant_uom_code"`
	CoolantUOMDescription string  `json:"coolant_uom_description"`

	Hydraulic               string  `json:"hydraulic"`
	HydraulicCapacity       float64 `json:"hydraulic_capacity"`
	IDHydraulicUOM          uint    `json:"id_hydraulic_uom"`
	HydraulicUOMCode        string  `json:"hydraulic_uom_code"`
	HydraulicUOMDescription string  `json:"hydraulic_uom_description"`

	DimensionFront               float64 `json:"dimension_front"`
	IDDimensionFrontUOM          uint    `json:"id_dimension_front_uom"`
	DimensionFrontUOMCode        string  `json:"dimension_front_uom_code"`
	DimensionFrontUOMDescription string  `json:"dimension_front_uom_description"`

	DimensionSide               float64 `json:"dimension_side"`
	IDDimensionSideUOM          uint    `json:"id_dimension_side_uom"`
	DimensionSideUOMCode        string  `json:"dimension_side_uom_code"`
	DimensionSideUOMDescription string  `json:"dimension_side_uom_description"`

	DetailCreatedAt     time.Time `json:"detail_created_at"`
	DetailUpdatedAt     time.Time `json:"detail_updated_at"`
	DetailCreatedby     uint      `json:"detail_createdby"`
	DetailCreatedbyName string    `json:"detail_createdby_name"`
	DetailUpdatedby     uint      `json:"detail_updatedby"`
	DetailUpdatedbyName string    `json:"detail_updatedby_name"`

	IDReason        uint      `json:"id_reason"`
	Remarks         string    `json:"remarks"`
	StatusCreatedAt time.Time `json:"status_created_at"`
	StatusUpdatedAt time.Time `json:"status_updated_at"`

	ReasonKey         string `json:"reason_key"`
	ReasonCode        string `json:"reason_code"`
	ReasonDescription string `json:"reason_description"`
	ReasonRemarks     string `json:"reason_remarks"`

	ApprovalID            uint      `json:"approval_id"`
	RefTable              string    `json:"ref_table"`
	RefID                 uint      `json:"ref_id"`
	ApprovalKey           string    `json:"approval_key"`
	ApprovalMessage       string    `json:"approval_message"`
	ApprovalCreatedby     uint      `json:"approval_createdby"`
	ApprovalCreatedbyName string    `json:"approval_createdby_name"`
	ApprovalCreatedAt     time.Time `json:"approval_created_at"`
	ApprovalStatus        string    `json:"approval_status"`
	ApprovalAction        string    `json:"approval_action"`
	ApprovalCount         int       `json:"approval_count"`
	ApprovalLevel         int       `json:"approval_level"`
}
