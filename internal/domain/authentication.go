package domain

import "time"

// Authentication represents refresh token stored in database
type Authentication struct {
	Token     string    `json:"token"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
