package model

type MstRolePermission struct {
	IDRole   uint `json:"id_role,omitempty"`
	IDMenu   uint `json:"id_menu,omitempty"`
	IsCreate bool `json:"is_create,omitempty"`
	IsUpdate bool `json:"is_update,omitempty"`
	IsDelete bool `json:"is_delete,omitempty"`

	Role *MstRole `gorm:"foreignKey:IDRole;references:ID" json:"role,omitempty"`
	Menu *MstMenu `gorm:"foreignKey:IDMenu;references:ID" json:"menu,omitempty"`
}
