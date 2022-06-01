package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/controller"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/middleware"
	"github.com/labstack/echo/v4"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Error loading .env file: %s\n", err.Error())
	}

	port := ":" + os.Getenv("PORT")

	registerController := controller.NewRegisterController()

	e := echo.New()

	if os.Getenv("ENV") == "production" {
		middleware.BodyLimit(e)
		middleware.Gzip(e)
		middleware.RateLimiter(e)
		middleware.Recover(e)
		middleware.Secure(e)
		middleware.RemoveTrailingSlash(e)
	} else {
		middleware.Logger(e)
		middleware.RemoveTrailingSlash(e)
	}

	g := e.Group("/api/v1")

	registerController.Route(g)

	e.Logger.Fatal(e.Start(port))
}
