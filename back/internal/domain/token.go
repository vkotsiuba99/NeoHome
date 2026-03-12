package domain

type TokenClaims struct {
	UserID    int64
	Role      string
	ExpiresAt int64
	IssuedAt  int64
}
