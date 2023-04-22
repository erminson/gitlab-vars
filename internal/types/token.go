package types

//import "time"

type Scope string

type Token struct {
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	Revoked    bool    `json:"revoked"`
	CreatedAt  Time    `json:"created_at"`
	Scopes     []Scope `json:"scopes"`
	UserId     int64   `json:"user_id"`
	LastUsedAt Time    `json:"last_used_at"`
	Active     bool    `json:"active"`
	ExpiresAt  string  `json:"expires_at"`
}
