package auth

import "context"

type TokenClaims struct {
	UserID   string
	Username string
}

type TokenIssuer interface {
	Issue(ctx context.Context, claims TokenClaims) (string, error)
}

type TokenParser interface {
	Parse(token string) (TokenClaims, error)
}
