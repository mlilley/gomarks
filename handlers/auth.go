package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/mlilley/gomarks/services"
	"net/http"
	"regexp"
)

func Token(authService services.AuthService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var input struct {
			Email    string `json:"email" form:"email"`
			Password string `json:"password" form:"password"`
			DeviceId string `json:"device_id" form:"device_id"`
		}

		err := c.Bind(&input)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		if input.Email == "" || input.Password == "" || input.DeviceId == "" || !isEmail(input.Email) || !isUUID(input.DeviceId) {
			return c.NoContent(http.StatusBadRequest)
		}

		_, _, accessToken, refreshToken, err := authService.Login(input.Email, input.Password, input.DeviceId)
		if err != nil {
			return c.NoContent(http.StatusUnauthorized)
		}

		return c.JSON(http.StatusOK, &map[string]string{
			"accessToken": accessToken,
			"refreshToken": refreshToken,
		})
	}
}

func Refresh(authService services.AuthService) echo.HandlerFunc {
	return func(c echo.Context) error {
		// obtain refreshToken from header/body?
		// ...

		user, deviceId, accessToken, authService.Refresh(refreshToken)
	}
}

var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var rxUUID = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")

func isEmail(s string) bool {
	return len(s) < 255 && rxEmail.MatchString(s)
}

func isUUID(s string) bool {
	return rxUUID.MatchString(s)
}
