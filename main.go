package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/controller"
	_ "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/docs"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/middleware"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title           Forum Group Discussion API
// @version         1.0
// @description     API for Forum Group Discussion
// @termsOfService  http://swagger.io/terms/

// @contact.name   Kelompok 22
// @contact.url    http://www.swagger.io/support
// @contact.email  erikriosetiawan15@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization

// @host      localhost:3000
// @BasePath  /api/v1
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

	e.GET("/*", echoSwagger.WrapHandler)

	g := e.Group("/api/v1")

	registerController.Route(g)

	e.Logger.Fatal(e.Start(port))
}
