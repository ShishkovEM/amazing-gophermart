package models

import (
	"time"

	"github.com/google/uuid"
)

type Withdraw struct {
	UserID   uuid.UUID `json:"user_id,omitempty" db:"user_id"`
	OrderNum string    `json:"order" db:"order_num"`
	Withdraw float32   `json:"sum" db:"withdrawal"`
	Created  time.Time `json:"uploaded_at" db:"withdrawn_at"`
}

type WithdrawDB struct {
	OrderNum string    `json:"order" db:"order_num"`
	Withdraw float32   `json:"sum" db:"withdrawal"`
	Created  time.Time `json:"processed_at" db:"withdrawn_at"`
}
