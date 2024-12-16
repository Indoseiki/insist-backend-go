package dto

type ApprovalUsers struct {
	IDApproval uint   `json:"id_approval"`
	IDUser     []uint `json:"id_user"`
}
