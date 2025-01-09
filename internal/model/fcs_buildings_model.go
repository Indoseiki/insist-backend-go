package model

type MstFCSBuilding struct {
	IDFCS      uint `json:"id_fcs,omitempty"`
	IDBuilding uint `json:"id_building,omitempty"`

	FCS      *MstFCS      `gorm:"foreignKey:IDFCS;references:ID" json:"fcs,omitempty"`
	Building *MstBuilding `gorm:"foreignKey:IDBuilding;references:ID" json:"building,omitempty"`
}
