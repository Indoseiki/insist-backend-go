package model

import (
	"time"
)

type MstEmployee struct {
	Number     string     `gorm:"primaryKey" json:"number"`
	Name       string     `json:"name,omitempty"`
	Division   string     `json:"division,omitempty"`
	Department string     `json:"department,omitempty"`
	Position   string     `json:"position,omitempty"`
	IsActive   bool       `json:"is_active,omitempty"`
	Service    string     `json:"service,omitempty"`
	Education  string     `json:"education,omitempty"`
	Birthday   time.Time  `json:"birthday,omitempty"`
	CreatedAt  *time.Time `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt  *time.Time `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
}
