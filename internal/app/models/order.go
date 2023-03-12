package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	UserID   uuid.UUID   `json:"user_id,omitempty" db:"user_id"`
	OrderNum string      `json:"number" db:"order_num"`
	Status   string      `json:"status" db:"status"`
	Accrual  interface{} `json:"accrual,omitempty" db:"accrual"`
	Created  time.Time   `json:"uploaded_at" db:"created_at"`
}

type OrderDB struct {
	OrderNum string      `json:"number" db:"order_num"`
	Status   string      `json:"status" db:"status"`
	Accrual  interface{} `json:"accrual,omitempty" db:"accrual"`
	Created  time.Time   `json:"uploaded_at" db:"created"`
}

type ProcessingOrder struct {
	OrderNum string      `json:"order"`
	Status   string      `json:"status"`
	Accrual  interface{} `json:"accrual,omitempty"`
}
