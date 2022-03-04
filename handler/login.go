package handler

import (
	"himatro-api/controller"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	NPM, plainPassword, err := controller.ExtractLoginPayload(c)

	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorMessage{
			OK:      false,
			Message: err.Error(),
		})
	}

	encryptedPassword, err := controller.GetUserPassword(NPM)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorMessage{
			OK:      false,
			Message: err.Error(),
		})
	}

	err = controller.ValidatePassword(plainPassword, encryptedPassword)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorMessage{
			OK:      false,
			Message: err.Error(),
		})
	}

	loginToken, err := controller.CreateLoginToken(c, NPM)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorMessage{
			OK:      false,
			Message: "Server error when creating login token.",
		})
	}

	return c.JSON(http.StatusOK, LoginTokenResp{
		OK:    true,
		Token: loginToken,
	})
}
