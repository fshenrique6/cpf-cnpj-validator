package models

import "time"

type Document struct {
	ID          int       `json:"id"`
	Number      string    `json:"number" binding:"required"`
	Type        string    `json:"type"`
	Blocklisted bool      `json:"blocklisted"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}