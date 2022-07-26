package auth

import (
	"context"

	"github.com/Himatro2021/API/internal/rbac"
)

type contextKey string

const (
	userCtxKey           contextKey = "github.com/Himatro2021/API/internal/auth.User"
	_headerAuthorization string     = "Authorization"
	_authScheme          string     = "Bearer"
)

// User define any data related for identifiying user
type User struct {
	ID    int64     `json:"id"`
	Email string    `json:"email"`
	Role  rbac.Role `json:"role"`
}

// HasAccess check if user has access to the given resource and action pair
func (u *User) HasAccess(resource rbac.Resource, action rbac.Action) bool {
	if string(u.Role) == "" {
		return false
	}

	return rbac.HasAccess(u.Role, resource, action)
}

// SetUserToCtx self explained
func SetUserToCtx(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, userCtxKey, user)
}

// GetUserFromCtx self explained
func GetUserFromCtx(ctx context.Context) *User {
	user, ok := ctx.Value(userCtxKey).(User)
	if !ok {
		return nil
	}

	return &user
}
