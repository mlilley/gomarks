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
	deviceRepo := repos.NewDeviceRepo(db)

	authService := services.NewAuthService(userRepo, deviceRepo, secret)
	//userService := services.NewUserService(userRepo)
	markService := services.NewMarkService(markRepo)
	authorizeMiddleware := handlers.AuthorizeMiddleware(authService)

	e.Use(middleware.Logger())

	//e.POST("/token", handlers.CreateToken(authService))

	marksGroup := e.Group("/marks", authorizeMiddleware)
	marksGroup.GET("", handlers.GetMarks(markService))
	marksGroup.GET("/:id", handlers.GetMark(markService)).Name = "GetMark"
	marksGroup.POST("", handlers.CreateMark(markService))
	marksGroup.PUT("/:id", handlers.UpdateMark(markService))
	marksGroup.DELETE("/:id", handlers.DeleteMark(markService))
	marksGroup.DELETE("/", handlers.DeleteMarks(markService))

	server := &http.Server{
		Addr:         "localhost:8080",
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
	}

	e.Logger.Fatal(e.StartServer(server))
}
