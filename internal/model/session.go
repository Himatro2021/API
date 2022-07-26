package model

import (
	"context"
	"time"
)

// Session represent what is stored in db as a session information
type Session struct {
	ID                    int64     `json:"-"`
	UserID                int64     `json:"-"`
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	AccessTokenExpiredAt  time.Time `json:"access_token_expired_at"`
	RefreshTokenExpiredAt time.Time `json:"refresh_token_expired_at"`
}

// IsAccessTokenExpired compare session access token against time.Now()
// return true if the token is expired by now
func (s *Session) IsAccessTokenExpired() bool {
	if s == nil {
		return true
	}

	now := time.Now()
	return now.After(s.AccessTokenExpiredAt)
}

// IsRefreshTokenExpired compare session refresh token against time.Now()
// return true if the token is expired by now
func (s *Session) IsRefreshTokenExpired() bool {
	if s == nil {
		return true
	}

	now := time.Now()
	return now.After(s.RefreshTokenExpiredAt)
}

// SessionRepository :nodoc:
type SessionRepository interface {
	Create(ctx context.Context, session *Session) error
	FindByAccessToken(ctx context.Context, token string) (*Session, error)
}
