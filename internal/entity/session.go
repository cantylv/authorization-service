package entity

import (
	"time"
)

type Session struct {
	ID            int
	UserID        string
	RefreshToken  string
	UserIpAddress string
	Fingerprint   string
	ExpiresAt     time.Time
	CreatedAt     time.Time
}
