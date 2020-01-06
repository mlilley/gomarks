package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/mlilley/gomarks/services"
	"net/http"
)

func HandleCreateToken(authService services.AuthService) echo.HandlerFunc {
	return func(c echo.Context) error {
		var input struct{
			Username string `json:"username" form:"username"`
			Password string `json:"password" form:"password"`
		}

		err := c.Bind(&input)
		if err != nil {
			return c.NoContent(http.StatusBadRequest)
		}

		token, err := authService.Authorize(input.Username, input.Password)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, &map[string]string{
				"error": "invalid credentials",
			})
		}

		// send token back in body
		return c.JSON(http.StatusOK, &map[string]string{"token": token})
	}
}