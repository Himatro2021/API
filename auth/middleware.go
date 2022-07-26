package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/Himatro2021/API/internal/model"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// Middleware :nodoc:
type Middleware struct {
	sessionRepo model.SessionRepository
	userRepo    model.UserRepository
}

// NewMiddleware return new instance of Middleware
func NewMiddleware(sessionRepo model.SessionRepository, userRepo model.UserRepository) *Middleware {
	return &Middleware{
		sessionRepo: sessionRepo,
		userRepo:    userRepo,
	}
}

// UserSessionMiddleware set user to context if access token valid. otherwise pass
func (am *Middleware) UserSessionMiddleware() echo.MiddlewareFunc {
	return func(hf echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := am.getTokenFromHeader(c.Request())
			if token == "" {
				return hf(c)
			}

			ctx := c.Request().Context()

			session, err := am.getSessionFromAccessToken(ctx, token)
			if err != nil {
				return hf(c)
			}

			// just pass if expired. The next middleware should block the request
			// if needed
			if session.IsAccessTokenExpired() {
				return hf(c)
			}

			user, err := am.userRepo.GetUserByID(ctx, session.UserID)
			if err != nil {
				return hf(c)
			}

			ctx = SetUserToCtx(ctx, User{
				ID:    user.ID,
				Email: user.Email,
				Role:  user.GetRole(),
			})

			c.SetRequest(c.Request().WithContext(ctx))

			return hf(c)
		}
	}
}

// RejectUnauthorizedRequest if no user in context, return unauthorized error. otherwise pass
func (am *Middleware) RejectUnauthorizedRequest(skipURL []string) echo.MiddlewareFunc {
	return func(hf echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			url := ctx.Request().URL
			for _, skip := range skipURL {
				if url.Path == skip {
					return hf(ctx)
				}
			}

			user := GetUserFromCtx(ctx.Request().Context())
			if user == nil {
				return ErrUnauthorized
			}

			return hf(ctx)
		}
	}
}

func (am *Middleware) getTokenFromHeader(req *http.Request) string {
	authHeader := strings.Split(req.Header.Get(_headerAuthorization), " ")

	if len(authHeader) != 2 || authHeader[0] != _authScheme {
		return ""
	}

	return strings.TrimSpace(authHeader[1])
}

func (am *Middleware) getSessionFromAccessToken(ctx context.Context, token string) (*model.Session, error) {
	logger := logrus.WithField("token", token)

	session, err := am.sessionRepo.FindByAccessToken(ctx, token)
	if err != nil {
		logger.Warn(err)
		return nil, err
	}

	return session, nil
}
