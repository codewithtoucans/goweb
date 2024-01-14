package context

import (
	"context"
	"github.com/codewithtoucans/goweb/models"
)

type key string

const _userKey key = "user"

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, _userKey, user)
}

func User(ctx context.Context) *models.User {
	value := ctx.Value(_userKey)
	user, ok := value.(*models.User)
	if !ok {
		return nil
	}
	return user
}
