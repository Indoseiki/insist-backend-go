package model

type MstUserRole struct {
	IDUser uint `json:"id_user,omitempty"`
	IDRole uint `json:"id_role,omitempty"`

	User *MstUser `gorm:"foreignKey:IDUser;references:ID" json:"user,omitempty"`
	Role *MstRole `gorm:"foreignKey:IDRole;references:ID" json:"role,omitempty"`
}
