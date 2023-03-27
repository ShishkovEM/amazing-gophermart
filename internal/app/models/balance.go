package models

type Balance struct {
	Current  float32 `json:"current"`
	Withdraw float32 `json:"withdrawn"`
}
