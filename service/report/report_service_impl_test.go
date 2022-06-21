package report

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/config"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository/report"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository/thread"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository/user"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator"
)

var db *sql.DB
var reportRepo report.ReportRepository
var threadRepo thread.ThreadRepository
var userRepo user.UserRepository
var idGenerator generator.IDGenerator

func init() {
	godotenv.Load("./../../.env.example")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var err error
	if db, err = config.NewPostgreSQLDatabase(); err != nil {
		panic(err)
	}

	reportRepo = report.NewReportRepositoryImpl(db)
	threadRepo = thread.NewThreadRepositoryImpl(db)
	userRepo = user.NewUserRepositoryImpl(db)
	idGenerator = generator.NewNanoidIDGenerator()
}

func TestGetAll(t *testing.T) {
	var service ReportService = NewReportServiceImpl(reportRepo, userRepo, threadRepo, idGenerator)

	accessorRole := "admin"
	status := "review"
	page := 1
	limit := 10

	if pagination, err := service.GetAll(context.Background(), accessorRole, status, uint(page), uint(limit)); err != nil {
		t.Fatalf("Error happened: %s", err)
	} else {
		t.Logf("Pagination: %+v", pagination)
	}
}
