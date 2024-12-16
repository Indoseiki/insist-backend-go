package model

type MstApprovalUser struct {
	IDApproval uint `json:"id_approval,omitempty"`
	IDUser     uint `json:"id_user,omitempty"`

	Approval *MstApproval `gorm:"foreignKey:IDApproval;references:ID" json:"approval,omitempty"`
	User     *MstUser     `gorm:"foreignKey:IDUser;references:ID" json:"user,omitempty"`
}
