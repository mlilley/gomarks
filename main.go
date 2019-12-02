package main

import (
	"github.com/labstack/echo/v4/middleware"
	"github.com/mlilley/gomarks/database"
	"github.com/mlilley/gomarks/mark"
	"github.com/mlilley/gomarks/user"
	"net/http"
	"time"

	echo "github.com/labstack/echo/v4"
)

type AppContext struct {
	echo.Context
	users *user.Repository
	marks *mark.Repository
	secret *[]byte
}

func main() {
	e := echo.New()

	secret := []byte("0DSG5J@nUoV1zk@b3%zW@z$47qO5LH$67RVGZBgCoVZ&HDph9G$NqPneg9ujCRO!q8xwH!Wes9mL!#IwkQu$7b$nbTE5#wUzKn%RvsdhqgXCP7UIV%*O@iiy*mgHl&X1")

	db, err := database.New()
	if err != nil {
		e.Logger.Fatal(err)
	}

	users, err := user.NewRepo(db)
	if err != nil {
		e.Logger.Fatal(err)
	}

	marks, err := mark.NewRepo(db)
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(&AppContext{
				Context: c,
				users: &users,
				marks: &marks,
				secret: &secret,
			})
		}
	})

	e.Use(middleware.Logger())

	e.POST("/token", authHandler)
	e.GET("/marks", getMarksHandler)
	e.GET("/marks/:id", getMarkHandler).Name = "GetMark"
	e.POST("/marks", createMarkHandler)
	e.PUT("/marks/:id", updateMarkHandler)
	e.DELETE("/marks/:id", deleteMarkHandler)
	e.DELETE("/marks", deleteMarksHandler)

	server := &http.Server{
		Addr:         "localhost:8080",
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
	}

	e.Logger.Fatal(e.StartServer(server))
}
