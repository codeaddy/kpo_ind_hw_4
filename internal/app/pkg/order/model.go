package order

import "time"

type Order struct {
	ID              int       `db:"id" json:"id"`
	UserID          int       `db:"user_id" json:"user_id"`
	Status          string    `db:"status" json:"status"`
	SpecialRequests string    `db:"special_requests" json:"special_requests"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}
