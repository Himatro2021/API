package auth

import (
	"errors"
	"fmt"
	"himatro-api/internal/config"
	"himatro-api/internal/util"
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
		util.LogErr("ERROR", "Server failed to create signed token string", err.Error())
		return "", errors.New("server failed to create update token")
	}

	return signedToken, nil
}

func ExtractJWTPayload(token string, claims *UpdateAbsentListClaims) error {
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSigningKey()), nil
	})

	if err != nil {
		util.LogErr("INFO", fmt.Sprintf("Invalid token used: %s", token), err.Error())
		return fmt.Errorf("invalid token: %s", err.Error())
	}

	return nil
}
