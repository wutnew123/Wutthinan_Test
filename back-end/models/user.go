package models

import (
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:100;not null"`
	Password  string    `json:"-" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RevokedToken struct {
	TokenHash string    `json:"-" gorm:"primaryKey;size:64"`
	ExpiresAt time.Time `json:"expires_at;index"`
	CreatedAt time.Time `json:"created_at"`
}

func (RevokedToken) TableName() string {
	return "revoked_tokens"
}