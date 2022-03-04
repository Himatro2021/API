package auth

import (
	"errors"
	"os"
	"strconv"
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
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_SIGNING_KEY")))

	if err != nil {
		return "", errors.New("server failed to create login token")
	}

	return signedToken, nil
}

func getTokenExpSec() int {
	tokenExpSec, err := strconv.Atoi(os.Getenv("LOGIN_TOKEN_EXP_SEC"))

	if err != nil {
		return 604800 // 7 days
	}

	return tokenExpSec
}

func createClaims(NPM string) jwtCustomClaims {
	claims := jwtCustomClaims{
		NPM,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(getTokenExpSec())).Unix(),
			Issuer:    NPM,
		},
	}

	return claims
}
