package dto

type ChangePassword struct {
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}
