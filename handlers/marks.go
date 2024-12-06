package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/mlilley/gomarks/app"
	"github.com/mlilley/gomarks/services"
	"net/http"
)

func GetMarks(marksService services.MarkService) echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO:
		user := &app.User{ID: "1", Email: "user@test.com", PasswordHash: "", Active: true}

		marks, err := marksService.GetMarksForUser(user.ID)
		if err != nil {
			return err
		}

		c.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, "*")
		return c.JSON(http.StatusOK, marks)
	}
}

func GetMark(markService services.MarkService) echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO:
		user := &app.User{ID: "1", Email: "user@test.com", PasswordHash: "", Active: true}

		markId := c.Param("id")
		mark, err := markService.GetMarkByIDForUser(markId, user.ID)
		if err != nil {
			return err
		}
		if mark == nil {
			return echo.NewHTTPError(http.StatusNotFound)
		}

		return c.JSON(http.StatusOK, mark)
	}
}


func CreateMark(markService services.MarkService) echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO:
		user := &app.User{ID: "1", Email: "user@test.com", PasswordHash: "", Active: true}

		mark := new(app.Mark)
		err := c.Bind(mark)
		if err != nil {
			return err
		}

		mark, err = markService.CreateMarkForUser(mark, user.ID)
		if err != nil {
			return err
		}

		c.Response().Header().Set(echo.HeaderLocation, c.Echo().Reverse("GetMark", mark.ID))
		return c.NoContent(http.StatusCreated)
	}
}

func UpdateMark(markService services.MarkService) echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO:
		user := &app.User{ID: "1", Email: "user@test.com", PasswordHash: "", Active: true}

		markId := c.Param("id")
		mark := new(app.Mark)
		err := c.Bind(mark)
		if err != nil {
			return err
		}

		mark.ID = markId
		mark, err = markService.UpdateMarkForUser(mark, user.ID)
		if err != nil {
			return err
		}

		if mark == nil {
			return c.NoContent(http.StatusNotFound)
		}

		// TODO: generate and return GET mark url in location header
		return c.JSON(http.StatusOK, mark)
	}
}

func DeleteMark(markService services.MarkService) echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO:
		user := &app.User{ID: "1", Email: "user@test.com", PasswordHash: "", Active: true}

		markId := c.Param("id")
		ok, err := markService.DeleteMarkByIDForUser(markId, user.ID)
		if err != nil {
			return err
		}

		if !ok {
			return c.NoContent(http.StatusNotFound)
		}

		return c.NoContent(http.StatusOK)
	}
}

func DeleteMarks(markService services.MarkService) echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO:
		user := &app.User{ID: "1", Email: "user@test.com", PasswordHash: "", Active: true}

		err := markService.DeleteMarksForUser(user.ID)
		if err != nil {
			return err
		}

		return c.NoContent(http.StatusOK)
	}
}