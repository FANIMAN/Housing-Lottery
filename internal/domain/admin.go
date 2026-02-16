package domain

import "time"

type Admin struct {
	ID           string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}
