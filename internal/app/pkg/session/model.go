package session

import "time"

type Session struct {
	ID           int       `db:"id" json:"id"`
	UserID       int       `db:"user_id" json:"user_id"`
	SessionToken string    `db:"session_token" json:"session_token"`
	ExpiresAt    time.Time `db:"expires_at" json:"expires_at"`
}
