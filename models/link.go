package models

import (
	"time"
	"gorm.io/gorm"
)

// Link struct
type Link struct {
	ID         string    `gorm:"primaryKey"`
	TargetURL  string
	Shortened  string    `gorm:"not null;unique"`
	Clicks     int
	MaxClicks  int
	ExpiresAt  time.Time
	CreatedAt  time.Time
	UserID     uint      `gorm:"not null"`

	ClicksData []LinkClick `gorm:"foreignKey:LinkID"` // One-to-many

	gorm.Model
}


// Link data analysis struct
type LinkClick struct {
	ID        uint      `gorm:"primaryKey"`
	LinkID    string    `gorm:"index"` // Still a string, matches Link.ID
	Time      time.Time
	IP        string
	Referrer  string
	Device    string
	Country   string

	Link Link `gorm:"foreignKey:LinkID;references:ID;constraint:OnDelete:CASCADE"`
}

