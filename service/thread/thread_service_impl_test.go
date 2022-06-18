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

	accessorUserID := "u-ZrxmQS"

	page := 1
	limit := 20
	query := "Read"

	if pagination, err := service.GetAll(context.Background(), accessorUserID, uint(page), uint(limit), query); err != nil {
		t.Fatalf("Error happened: %s", err)
	} else {
		t.Logf("Pagination: %+v", pagination)
	}
}

func TestCreate(t *testing.T) {
	var service ThreadService = NewThreadServiceImpl(threadRepo, categoryRepo, idGenerator)

	accessorUserID := "u-ZrxmQS"

	p := payload.CreateThread{
		Title:       "Go Programming Language Going Hype",
		Description: "Currently Go Programming Language going hype, because it features.",
		CategoryID:  "c-abc",
	}

	if id, err := service.Create(context.Background(), accessorUserID, p); err != nil {
		t.Fatalf("Error happened: %s", err)
	} else {
		t.Logf("Successfully create a thread with ID %s", id)
	}
}

func TestGetByID(t *testing.T) {
	var service ThreadService = NewThreadServiceImpl(threadRepo, categoryRepo, idGenerator)

	accessorUserID := "u-ZrxmQS"

	id := "t-kxdoB7i"

	if thread, err := service.GetByID(context.Background(), accessorUserID, id); err != nil {
		t.Fatalf("Error happened: %s", err)
	} else {
		t.Logf("Thread: %+v", thread)
	}
}

func TestUpdate(t *testing.T) {
	var service ThreadService = NewThreadServiceImpl(threadRepo, categoryRepo, idGenerator)

	accessorUserID := "u-ZrxmQS"

	id := "t-kxdoB7i"

	p := payload.UpdateThread{
		Title:       "Go Programming Going Hype",
		Description: "Currently Go Programming going hype, because it features.",
		CategoryID:  "c-def",
	}

	if err := service.Update(context.Background(), accessorUserID, id, p); err != nil {
		t.Fatalf("Error happened: %s", err)
	} else {
		t.Logf("Successfully updated a thread with ID %s", id)
	}
}

func TestDelete(t *testing.T) {
	var service ThreadService = NewThreadServiceImpl(threadRepo, categoryRepo, idGenerator)

	accessorUserID := "u-ZrxmQS"
	role := "user"

	id := "t-45uolpR"

	if err := service.Delete(context.Background(), accessorUserID, role, id); err != nil {
		t.Fatalf("Error happened: %s", err)
	} else {
		t.Logf("Successfully deleted a thread with ID %s", id)
	}
}

func TestGetComments(t *testing.T) {
	var service ThreadService = NewThreadServiceImpl(threadRepo, categoryRepo, idGenerator)

	threadID := "t-abcdefg"

	if pagination, err := service.GetComments(context.Background(), threadID, 0, 0); err != nil {
		t.Fatalf("Error happened: %s", err)
	} else {
		t.Logf("Pagination: %+v", pagination)
	}
}

func TestChangeFollowingState(t *testing.T) {
	var service ThreadService = NewThreadServiceImpl(threadRepo, categoryRepo, idGenerator)

	threadID := "t-abcdefg"
	accessorUserID := "u-kt56R1"

	if err := service.ChangeFollowingState(context.Background(), threadID, accessorUserID); err != nil {
		t.Fatalf("Error happened: %s", err)
	} else {
		t.Log("Successfully change the following state")
	}
}
