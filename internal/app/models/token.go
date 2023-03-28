package models

type Token struct {
	TokenType   string `json:"token_type"`
	AuthToken   string `json:"auth_token"`
	GeneratedAt string `json:"generated_at"`
	ExpiresAt   string `json:"expires_at"`
}
