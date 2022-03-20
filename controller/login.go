package controller

import (
	"errors"
	"himatro-api/auth"
	"himatro-api/db"
	"himatro-api/models"

	_ "github.com/joho/godotenv/autoload"

	"github.com/labstack/echo/v4"
)

type LoginPayload struct {
	NPM      string `json:"NPM"`
	Password string `json:"password"`
}

func GetUserPassword(NPM string) (string, error) {
	user := models.User{}

	err := db.DB.Where("npm = ?", NPM).First(&user)

	if err.Error != nil {
		return "", errors.New("login credentials in not valid")
	}

	return user.Password, nil
}

func ExtractLoginPayload(c echo.Context) (string, string, error) {
	payload := new(LoginPayload)

	if err := c.Bind(payload); err != nil {
		return "", "", errors.New("payload is incorrect")
	}

	if payload.NPM == "" || payload.Password == "" {
		return "", "", errors.New("NPM and password must be supplied")
	}

	return payload.NPM, payload.Password, nil
}

func ValidatePassword(plain string, encrypted string) error {
	decrypted, err := auth.Decrypt(encrypted)

	if err != nil {
		return errors.New("credentials invalid")
	}

	if decrypted != plain {
		return errors.New("credentials invalid")
	}

	return nil
}

func CreateLoginToken(c echo.Context, NPM string) (string, error) {
	loginToken, err := auth.CreateLoginToken(NPM)

	if err != nil {
		return "", err
	}

	return loginToken, nil
}
