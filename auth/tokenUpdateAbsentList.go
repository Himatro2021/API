package auth

import (
	"errors"
	"fmt"
	"os"
	"strconv"
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
			ExpiresAt: time.Now().Add(time.Second * time.Duration(getIntegerEnvVar("UPDATE_ABSENT_LIST_TOKEN_EXP_SEC", 3600))).Unix(),
			Issuer:    NPM,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_SIGNING_KEY")))

	if err != nil {
		return "", errors.New("server failed to create update token")
	}

	return signedToken, nil
}

func ExtractJWTPayload(token string, claims *UpdateAbsentListClaims) error {
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_SIGNING_KEY")), nil
	})

	if err != nil {
		return fmt.Errorf("invalid token: %s", err.Error())
	}

	return nil
}

func getIntegerEnvVar(name string, defaultVal int) int {
	envVal, err := strconv.Atoi(os.Getenv(name))

	if err != nil {
		return defaultVal
	}

	return envVal
}
