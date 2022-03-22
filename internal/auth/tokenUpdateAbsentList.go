package auth

import (
	"errors"
	"fmt"
	"himatro-api/internal/config"
	"time"

	"github.com/golang-jwt/jwt"
)

type UpdateAbsentListClaims struct {
	NPM      string `json:"npm"`
	AbsentID uint   `json:"absentID"`
	jwt.StandardClaims
}

func CreateUpdateAbsentListToken(absentID int, NPM string) (string, error) {
	AbsentID := uint(absentID)
	claims := UpdateAbsentListClaims{
		NPM,
		AbsentID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(config.UpdateAbsentListTokenExpSec())).Unix(),
			Issuer:    NPM,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(config.JWTSigningKey()))

	if err != nil {
		return "", errors.New("server failed to create update token")
	}

	return signedToken, nil
}

func ExtractJWTPayload(token string, claims *UpdateAbsentListClaims) error {
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSigningKey()), nil
	})

	if err != nil {
		return fmt.Errorf("invalid token: %s", err.Error())
	}

	return nil
}
