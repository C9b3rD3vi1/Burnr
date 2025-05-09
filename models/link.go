package models

import (
	"time"

	//"gorm.io/gorm"
)

// Link struct
type Link struct {
	ID        string       `json:"id" gorm:"primaryKey"`
	Shortened string `json:"shortened" gorm:"unique;not null"`
	TargetURL  string `json:"original" gorm:"not null"`
	Clicks    int    `json:"clicks" gorm:"default:0"`
	MaxClicks int    `json:"max_clicks" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
	ExpireAt  time.Time `json:"expire_at" gorm:"index"`
}

