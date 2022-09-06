package models

import "time"

type PasswordReset struct {
	Email     string `json:"email"`
	Token     string `json:"token"`
	CreatedAt time.Time
}
