package thread

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/config"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository/thread"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator"
)

var db *sql.DB
var threadRepo thread.ThreadRepository
var idGenerator generator.IDGenerator

func init() {
	godotenv.Load("./../../.env.example")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var err error
	if db, err = config.NewPostgreSQLDatabase(); err != nil {
		panic(err)
	}

	threadRepo = thread.NewThreadRepositoryImpl(db)
	idGenerator = generator.NewNanoidIDGenerator()
}

func TestGetAll(t *testing.T) {
	var service ThreadService = NewThreadServiceImpl(threadRepo, idGenerator)

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
