package dto

type UserRoles struct {
	IDUser uint   `json:"id_user"`
	IDRole []uint `json:"id_role"`
}
