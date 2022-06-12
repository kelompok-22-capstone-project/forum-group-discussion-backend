package thread

import (
	"context"
	"database/sql"
	"testing"

	"github.com/joho/godotenv"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/config"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/entity"
)

var db *sql.DB

func init() {
	err := godotenv.Load("./../../.env.example")
	if err != nil {
		panic(err)
	}

	db, err = config.NewPostgreSQLDatabase()
	if err != nil {
		panic(err)
	}
}

func TestInsert(t *testing.T) {
	var repo ThreadRepository = NewThreadRepositoryImpl(db)

	thread := entity.Thread{
		ID:          "t-abcdefg",
		Title:       "Thread Title",
		Description: "This is thread description.",
		Creator: entity.User{
			ID: "u-ZrxmQS",
		},
		Category: entity.Category{
			ID: "c-abc",
		},
	}

	if err := repo.Insert(context.Background(), thread); err != nil {
		t.Fatalf("Error happened: %+v", err)
	} else {
		t.Logf("Thread with id %s successfully inserted", thread.ID)
	}
}
