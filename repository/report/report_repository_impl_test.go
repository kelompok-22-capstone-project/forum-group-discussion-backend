package report

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/config"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/entity"
)

var db *sql.DB

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	err := godotenv.Load("./../../.env.example")
	if err != nil {
		panic(err)
	}

	db, err = config.NewPostgreSQLDatabase()
	if err != nil {
		panic(err)
	}
}

func TestGetReportsWithPagination(t *testing.T) {
	var repo ReportRepository = NewReportRepositoryImpl(db)

	pageInfo := entity.PageInfo{Page: 1, Limit: 10}
	reportStatus := entity.Review

	if pagination, err := repo.GetReportsWithPagination(context.Background(), pageInfo, reportStatus); err != nil {
		t.Fatalf("Error happened: %+v", err)
	} else {
		t.Logf("Pagination: %+v", pagination)
	}
}

func TestInsert(t *testing.T) {
	var repo ReportRepository = NewReportRepositoryImpl(db)

	ID := "b-abcdefg"
	moderatorID := "m-1234"
	userID := "u-NplCrv"
	reason := "Bocil ini sangat meresahkan, tolong di banned min."

	if err := repo.Insert(context.Background(), ID, moderatorID, userID, reason); err != nil {
		t.Fatalf("Error happened: %+v", err)
	} else {
		t.Logf("Successfully inserted report with ID %s", ID)
	}
}
