package controller

import (
	"errors"
	"fmt"
	"himatro-api/internal/auth"
	"himatro-api/internal/contract"
	"himatro-api/internal/db"
	"himatro-api/internal/models"
	"himatro-api/internal/util"

	"github.com/labstack/echo/v4"
)

func GetUserPassword(NPM string) (string, error) {
	user := models.User{}

	err := db.DB.Where("npm = ?", NPM).First(&user)

	if err.Error != nil {
		util.LogErr("WARN", "Invalid login credentials were used", NPM)
		return "", errors.New("login credentials in not valid")
	}

	return user.Password, nil
}

func ExtractLoginPayload(c echo.Context) (string, string, error) {
	payload := new(contract.LoginPayload)

	if err := c.Bind(payload); err != nil {
		util.LogErr("WARN", "Invalid login payload were used", err.Error())
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
		util.LogErr("WARN", "Decryption process failed", err.Error())
		return errors.New("credentials invalid")
	}

	if decrypted != plain {
		util.LogErr("WARN", "Invalid login credentials were used", "")
		return errors.New("credentials invalid")
	}

	return nil
}

func CreateLoginToken(c echo.Context, NPM string) (string, error) {
	loginToken, err := auth.CreateLoginToken(NPM)

	if err != nil {
		util.LogErr("ERROR", fmt.Sprintf("Failed to create login token for NPM: %s", NPM), err.Error())
		return "", err
	}

	return loginToken, nil
}
