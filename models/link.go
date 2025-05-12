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

	// Optional reverse relationship
	ClicksData []LinkClick `gorm:"foreignKey:LinkID"`

	gorm.Model
}


// Link data analysis struct
type LinkClick struct {
	ID        uint      `gorm:"primaryKey"`
	LinkID    string    // This links back to Link.ID (which is a string)
	Time      time.Time
	IP        string
	Referrer  string
	Device    string
	Country   string

	// Optional: Add relation back to Link
	Link      Link      `gorm:"foreignKey:LinkID;constraint:OnDelete:CASCADE"`
}

