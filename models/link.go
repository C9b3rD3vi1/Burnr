package models

import (
	"time"
	"gorm.io/gorm"
)


// Link struct
type Link struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
	Shortened string `json:"shortened" gorm:"unique;not null"`
	Original  string `json:"original" gorm:"not null"`
	Clicks    int    `json:"clicks" gorm:"default:0"`
}