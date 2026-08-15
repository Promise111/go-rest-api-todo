package models

import "time"

type Users struct {
	ID        string    `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"password" db:"password"`
	Username  string    `json:"username" db:"username"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
