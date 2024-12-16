package model

type MstRoleMenu struct {
	IDRole uint `json:"id_role,omitempty"`
	IDMenu uint `json:"id_menu,omitempty"`

	Role *MstRole `gorm:"foreignKey:IDRole;references:ID" json:"role,omitempty"`
	Menu *MstMenu `gorm:"foreignKey:IDMenu;references:ID" json:"menu,omitempty"`
}
