package models

type Balance struct {
	Current  float32 `json:"current" db:"current"`
	Withdraw float32 `json:"withdrawn" db:"withdraw"`
}
