package main

import (
	"github.com/labstack/echo/v4"
	"github.com/mlilley/gomarks/auth"
	"net/http"
)

type inputData struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type outputData struct {
	Token string `json:"token"'`
}

func authHandler(c echo.Context) error {
	cc := c.(*AppContext)
	input := new(inputData)

	err := cc.Bind(input)
	if err != nil {
		return err
	}

	user, err := (*cc.users).FindByEmail(input.Username)
	if err != nil {
		return err
	}

	if user == nil || !user.Active {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	isAuthenticated := auth.CheckPassword(input.Password, user.PasswordHash)
	if !isAuthenticated {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	token, err := auth.GenerateToken(cc.secret, user.Email)
	if err != nil {
		return err
	}

	// send token back in Authorization header?
	//c.Response().Header().Set(echo.HeaderAuthorization, token)
	//return c.NoContent(http.StatusOK)

	// send token back in body
	return c.JSON(http.StatusOK, &outputData{Token: token})
}