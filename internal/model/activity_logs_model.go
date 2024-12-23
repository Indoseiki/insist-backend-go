package model

import "time"

type ActivityLog struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	IDUser    uint       `json:"id_user,omitempty"`
	IPAddress string     `json:"ip_address,omitempty"`
	Action    string     `json:"action,omitempty"`
	IsSuccess bool       `json:"is_success,omitempty"`
	Message   string     `json:"message,omitempty"`
	UserAgent string     `json:"user_agent,omitempty"`
	OS        string     `json:"os,omitempty"`
	CreatedAt *time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`

	User *MstUser `gorm:"foreignKey:IDUser;references:ID" json:"user,omitempty"`
}
