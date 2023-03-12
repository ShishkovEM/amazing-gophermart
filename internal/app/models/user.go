package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `json:"id" db:"id"`
	Username      string    `json:"login" db:"username"`
	Password      string    `json:"password" db:"pass"`
	Cookie        string    `json:"cookie" db:"cookie"`
	CookieExpires time.Time `json:"cookie_expires" db:"cookie_expires"`
}
