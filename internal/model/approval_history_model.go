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
