package auth

import "context"

type ctxKey struct{}

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, ctxKey{}, userID)
}

func UserID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(ctxKey{}).(string)
	return id, ok && id != ""
}

func RequireUserID(ctx context.Context) (string, error) {
	id, ok := UserID(ctx)
	if !ok {
		return "", ErrUnauthorized
	}
	return id, nil
}
