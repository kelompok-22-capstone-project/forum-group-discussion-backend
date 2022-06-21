package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/config"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/controller"
	_ "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/docs"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/middleware"
	ar "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository/admin"
	cr "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository/category"
	rr "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository/report"
	tr "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository/thread"
	ur "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository/user"
	as "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service/admin"
	cs "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service/category"
	rs "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service/report"
	ts "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service/thread"
	us "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service/user"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator"
	_ "github.com/kelompok-22-capstone-project/forum-group-discussion-backend/validation"
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

// @securityDefinitions.apikey  ApiKey
// @in                          header
// @name                        API-Key

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization

// @host      moot-rest-api.herokuapp.com
// @BasePath  /api/v1
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Error loading .env file: %s\n", err.Error())
	}

	db, err := config.NewPostgreSQLDatabase()
	if err != nil {
		log.Fatalln(err.Error())
	} else {
		log.Printf("Successfully connected to database with instance address: %p", db)
	}

	port := ":" + os.Getenv("PORT")

	idGenerator := generator.NewNanoidIDGenerator()
	passwordGenerator := generator.NewBcryptPasswordGenerator()
	tokenGenerator := generator.NewJWTTokenGenerator()

	userRepository := ur.NewUserRepositoryImpl(db)
	categoryRepository := cr.NewCategoryRepositoryImpl(db)
	threadRepository := tr.NewThreadRepositoryImpl(db)
	reportRepository := rr.NewReportRepositoryImpl(db)
	adminRepository := ar.NewAdminRepositoryImpl(db)

	userService := us.NewUserServiceImpl(userRepository, threadRepository, idGenerator, passwordGenerator, tokenGenerator)
	categoryService := cs.NewCategoryServiceImpl(categoryRepository, threadRepository, idGenerator)
	threadService := ts.NewThreadServiceImpl(threadRepository, categoryRepository, userRepository, idGenerator)
	reportService := rs.NewReportServiceImpl(reportRepository, userRepository, threadRepository, idGenerator)
	adminService := as.NewAdminServiceImpl(adminRepository)

	registerController := controller.NewRegisterController(userService)
	loginController := controller.NewLoginController(userService)
	usersController := controller.NewUsersController(userService, tokenGenerator)
	categoriesController := controller.NewCategoriesController(categoryService, tokenGenerator)
	threadsController := controller.NewThreadsController(threadService, tokenGenerator)
	adminController := controller.NewAdminController(adminService, tokenGenerator)
	reportsController := controller.NewReportsController(reportService, tokenGenerator)
	guestController := controller.NewGuestController(threadService, userService)

	e := echo.New()

	if os.Getenv("ENV") == "production" {
		middleware.CORS(e)
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

	g := e.Group("/api/v1", middleware.KeyAuth())

	registerController.Route(g)
	loginController.Route(g)
	usersController.Route(g)
	categoriesController.Route(g)
	threadsController.Route(g)
	adminController.Route(g)
	reportsController.Route(g)
	guestController.Route(g)

	e.Logger.Fatal(e.Start(port))
}
