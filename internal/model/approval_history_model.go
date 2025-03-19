package model

import (
	"time"
)

type ApprovalHistory struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	IDApproval  uint       `gorm:"not null" json:"id_approval"`
	RefTable    string     `gorm:"not null" json:"ref_table"`
	RefID       uint       `gorm:"not null" json:"ref_id"`
	Key         string     `gorm:"not null" json:"key"`
	Message     string     `gorm:"not null" json:"message"`
	IDCreatedby uint       `json:"id_createdby"`
	CreatedAt   *time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`

	Approval  *MstApproval `gorm:"foreignKey:ID;references:IDApproval" json:"approval,omitempty"`
	CreatedBy *MstUser     `gorm:"foreignKey:ID;references:IDCreatedby" json:"created_by,omitempty"`
}

type ViewApprovalNotification struct {
	ID                  uint    `json:"id"`
	MenuName            string  `json:"menu_name"`
	MenuPath            string  `json:"menu_path"`
	RefTable            string  `json:"ref_table"`
	RefID               uint    `json:"ref_id"`
	Key                 string  `json:"key"`
	Message             string  `json:"message"`
	ApprovalID          uint    `json:"approval_id"`
	CurrentLevel        int     `json:"current_level"`
	Status              string  `json:"status"`
	Action              string  `json:"action"`
	CurrentApproverID   *uint   `json:"current_approver_id,omitempty"`
	CurrentApproverName *string `json:"current_approver_name,omitempty"`
	NextApprovalID      *uint   `json:"next_approval_id,omitempty"`
	NextLevel           *int    `json:"next_level,omitempty"`
	NextIDUser          *uint   `json:"next_id_user,omitempty"`
	NextApprovalName    *string `json:"next_approval_name,omitempty"`
	NextAction          *string `json:"next_action,omitempty"`
	NextStatus          *string `json:"next_status,omitempty"`
}
