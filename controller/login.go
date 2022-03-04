package controller

import (
	"errors"
	"himatro-api/auth"
	"himatro-api/db"
	"himatro-api/models"

	_ "github.com/joho/godotenv/autoload"

	"github.com/labstack/echo/v4"
)

func GetUserPassword(NPM string) (string, error) {
	user := models.User{}

	err := db.DB.Where("npm = ?", NPM).First(&user)

	if err.Error != nil {
		return "", errors.New("login credentials in not valid")
	}

	return user.Password, nil
}

func ExtractLoginPayload(c echo.Context) (string, string, error) {
	NPM := c.FormValue("NPM")
	password := c.FormValue("password")

	if NPM == "" || password == "" {
		return NPM, password, errors.New("NPM and password must be supplied")
	}

	return NPM, password, nil
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
