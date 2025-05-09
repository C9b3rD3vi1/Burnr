package models

import (
	"time"

	//"gorm.io/gorm"
)

// Link struct
type Link struct {
	ID        string    `gorm:"primaryKey"`
	TargetURL string
	Clicks    int
	MaxClicks int
	ExpiresAt time.Time
	CreatedAt time.Time
}
