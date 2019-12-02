package main

import (
	"github.com/labstack/echo/v4"
	mark "github.com/mlilley/gomarks/mark"
	"net/http"
)

type createMarkInputData struct {
	Title string `json:"title"`
	URL string `json:"url"`
}

type updateMarkInputData struct {
	ID string `json:"id"`
	Title string `json:"title"`
	URL string `json:"url"`
}

func getMarksHandler(c echo.Context) error {
	cc := c.(*AppContext)

	marks, err := (*cc.marks).FindAll()
	if err != nil {
		return err
	}

	return cc.JSON(http.StatusOK, marks)
}

func getMarkHandler(c echo.Context) error {
	cc := c.(*AppContext)
	id := cc.Param("id")

	m, err := (*cc.marks).FindByID(id)
	if err != nil {
		return err
	}

	if m == nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	return cc.JSON(http.StatusOK, m)
}

func createMarkHandler(c echo.Context) error {
	cc := c.(*AppContext)

	input := new(createMarkInputData)
	err := cc.Bind(input)
	if err != nil {
		return err // TODO: probably not an InternalServer error
	}

	m := mark.Mark{
		Title: input.Title,
		URL: input.URL,
	}

	mm, err := (*cc.marks).Create(&m)
	if err != nil {
		return err // TODO probably not an InternalServer error
	}

	url := cc.Echo().Reverse("GetMark", mm.ID)

	c.Response().Header().Set(echo.HeaderLocation, url)
	return cc.NoContent(http.StatusCreated)
}

func updateMarkHandler(c echo.Context) error {
	cc := c.(*AppContext)
	id := cc.Param("id")

	input := new(updateMarkInputData)
	err := cc.Bind(input)
	if err != nil {
		return err // TODO
	}

	if id != input.ID {
		return echo.NewHTTPError(http.StatusBadRequest, "Cannot change id")
	}

	m := mark.Mark{
		ID: input.ID,
		Title: input.Title,
		URL: input.URL,
	}

	err = (*cc.marks).Update(&m)
	if err != nil {
		return err // TODO
	}

	return cc.NoContent(http.StatusOK)
}

func deleteMarkHandler(c echo.Context) error {
	cc := c.(*AppContext)
	id := cc.Param("id")

	exists, err := (*cc.marks).DeleteByID(id)
	if err != nil {
		return err
	}

	if !exists {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	return cc.NoContent(http.StatusOK)
}

func deleteMarksHandler(c echo.Context) error {
	cc := c.(*AppContext)

	err := (*cc.marks).DeleteAll()
	if err != nil {
		return err
	}

	return cc.NoContent(http.StatusOK)
}
