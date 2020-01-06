package main

import (
	"database/sql"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mlilley/gomarks/handlers"
	"github.com/mlilley/gomarks/repos"
	"github.com/mlilley/gomarks/services"
	"net/http"
	"time"
	"github.com/joho/godotenv"
	"os"
	"github.com/labstack/echo/v4"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	e := echo.New()

	err := godotenv.Load()
	if err != nil {
		e.Logger.Fatal(err)
	}

	secretEnv, ok := os.LookupEnv("GOMARKS_SECRET")
	if !ok {
		e.Logger.Fatal("GOMARKS_SECRET not set")
	}
	secret := []byte(secretEnv)

	db, err := sql.Open("sqlite3", "gomarks.db")
	if err != nil {
		e.Logger.Fatal(err)
	}

	userRepo := repos.NewUserRepo(db)
	markRepo := repos.NewMarkRepo(db)
	authService := services.NewAuthService(userRepo, secret)
	//userService := services.NewUserService(userRepo)
	markService := services.NewMarkService(markRepo)

	e.Use(middleware.Logger())

	e.POST("/token", handlers.HandleCreateToken(authService))
	e.GET("/marks", handlers.HandleGetMarks(markService))
	//e.GET("/marks/:id", h.GetMark).Name = "GetMark"
	//e.POST("/marks", h.CreateMark)
	//e.PUT("/marks/:id", h.UpdateMark)
	//e.DELETE("/marks/:id", h.DeleteMark)
	//e.DELETE("/marks", h.DeleteMarks)

	server := &http.Server{
		Addr:         "localhost:8080",
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
	}

	e.Logger.Fatal(e.StartServer(server))
}
