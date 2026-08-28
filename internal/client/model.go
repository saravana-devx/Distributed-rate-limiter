package client

import (
	"gorm.io/gorm"
)

type algorithmType string

const (
	FixedWindow   algorithmType = "fixed_window"
	SlidingWindow algorithmType = "sliding_window"
	TokenBucket   algorithmType = "token_bucket"
)

type Client struct {
	ID            int            `gorm:"primaryKey;autoIncrement;"`
	ClientID      string         `gorm:"type:varchar(100);not null;unique"`
	Algorithm     algorithmType  `gorm:"type:varchar(20);not null;default:'fixed_window'"`
	Limit         int            `gorm:"not null;default:10"`
	WindowSeconds int            `gorm:"not null;default:60"`
	APIKey        string         `gorm:"type:varchar(255);not null;unique"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}
