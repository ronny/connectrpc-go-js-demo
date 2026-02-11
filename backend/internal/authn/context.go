package authn

import (
	"context"
	"errors"

	"example.com/internal/service"
)

type ctxKeyType int

const (
	userKey ctxKeyType = iota
)

var ErrUnauthenticated = errors.New("unauthenticated")

func NewContext(ctx context.Context, user *service.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func RequireUser(ctx context.Context) (*service.User, error) {
	user, ok := ctx.Value(userKey).(*service.User)
	if !ok || user == nil || user.Email == "" {
		return nil, ErrUnauthenticated
	}

	return user, nil
}
