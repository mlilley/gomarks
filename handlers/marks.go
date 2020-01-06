package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/mlilley/gomarks/app"
	"github.com/mlilley/gomarks/services"
	"net/http"
)

func HandleGetMarks(marksService services.MarkService) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := &app.User{ID: "1", Email: "user@test.com", PasswordHash: "", Active: true}

		marks, err := marksService.GetMarksForUser(user)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, &map[string][]app.Mark{"marks":marks})
	}
}