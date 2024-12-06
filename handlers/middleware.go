package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/mlilley/gomarks/services"
	"strings"
)

func AuthorizeMiddleware(authService services.AuthService) func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// disabling auth for the time being...
			if (true) {
				return next(c)
			}
			// end disable

			token := parseToken(c.Request().Header.Get(echo.HeaderAuthorization))
			if token == "" {
				return echo.ErrUnauthorized
			}
			user, _, err := authService.Authorize(token)
			if err != nil {
				return echo.ErrUnauthorized
			}
			c.Set("user", user)
			return next(c)
		}
	}
}

func parseToken(header string) string {
	if strings.HasPrefix(header, "Bearer ") {
		return strings.TrimPrefix(header, "Bearer ")
	}
	return ""
}