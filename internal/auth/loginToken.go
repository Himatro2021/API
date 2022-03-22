package auth

import (
	"errors"
	"himatro-api/internal/config"
	"time"

	"github.com/golang-jwt/jwt"
)

type jwtCustomClaims struct {
	NPM string `json:"npm"`
	jwt.StandardClaims
}

func CreateLoginToken(NPM string) (string, error) {
	claims := createClaims(NPM)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(config.JWTSigningKey()))

	if err != nil {
		return "", errors.New("server failed to create login token")
	}

	return signedToken, nil
}

func createClaims(NPM string) jwtCustomClaims {
	claims := jwtCustomClaims{
		NPM,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(config.LoginTokenExpSec())).Unix(),
			Issuer:    NPM,
		},
	}

	return claims
}
