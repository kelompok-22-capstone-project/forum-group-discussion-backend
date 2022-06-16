package thread

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/config"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/payload"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository/category"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository/thread"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator"
)

var db *sql.DB
var threadRepo thread.ThreadRepository
var categoryRepo category.CategoryRepository
var idGenerator generator.IDGenerator

func init() {
	godotenv.Load("./../../.env.example")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var err error
	if db, err = config.NewPostgreSQLDatabase(); err != nil {
		panic(err)
	}

	threadRepo = thread.NewThreadRepositoryImpl(db)
	categoryRepo = category.NewCategoryRepositoryImpl(db)
	idGenerator = generator.NewNanoidIDGenerator()
}

func TestGetAll(t *testing.T) {
	var service ThreadService = NewThreadServiceImpl(threadRepo, categoryRepo, idGenerator)

	tp := generator.TokenPayload{
		ID:       "u-ZrxmQS",
		Username: "erikrios",
		Role:     "user",
		IsActive: true,
	}

	page := 1
	limit := 20
	query := "Read"

	if pagination, err := service.GetAll(context.Background(), tp, uint(page), uint(limit), query); err != nil {
		t.Fatalf("Error happened: %s", err)
	} else {
		t.Logf("Pagination: %+v", pagination)
	}
}

func TestCreate(t *testing.T) {
	var service ThreadService = NewThreadServiceImpl(threadRepo, categoryRepo, idGenerator)

	tp := generator.TokenPayload{
		ID:       "u-ZrxmQS",
		Username: "erikrios",
		Role:     "user",
		IsActive: true,
	}

	p := payload.CreateThread{
		Title:       "Go Programming Language Going Hype",
		Description: "Currently Go Programming Language going hype, because it features.",
		CategoryID:  "c-abc",
	}

	if id, err := service.Create(context.Background(), tp, p); err != nil {
		t.Fatalf("Error happened: %s", err)
	} else {
		t.Logf("Successfully create a thread with ID %s", id)
	}
}
