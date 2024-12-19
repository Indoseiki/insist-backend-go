package dto

type MenuWithPermissions struct {
	Label    string `json:"label"`
	Path     string `json:"path"`
	IDRole   uint   `json:"id_role"`
	IDMenu   uint   `json:"id_menu"`
	IsCreate bool   `json:"is_create"`
	IsUpdate bool   `json:"is_update"`
	IsDelete bool   `json:"is_delete"`
}
