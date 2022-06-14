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

func TestFindAllWithPagination(t *testing.T) {
	var repo ThreadRepository = NewThreadRepositoryImpl(db)

	pageInfo := entity.PageInfo{
		Limit: 10,
		Page:  1,
	}

	if pagination, err := repo.FindAllWithPagination(context.Background(), "u-ZrxmQS", pageInfo); err != nil {
		t.Fatalf("Error happened: %+v", err)
	} else {
		t.Logf("Pagination: %+v", pagination)
	}
}

func TestFindByID(t *testing.T) {
	var repo ThreadRepository = NewThreadRepositoryImpl(db)

	if thread, err := repo.FindByID(context.Background(), "u-ZrxmQS", "t-abcdefg"); err != nil {
		t.Fatalf("Error happened: %+v", err)
	} else {
		t.Logf("Thread: %+v", thread)
	}
}

func TestUpdate(t *testing.T) {
	var repo ThreadRepository = NewThreadRepositoryImpl(db)

	threadID := "t-abcdefg"

	thread := entity.Thread{}
	thread.Title = "Go Programming Become Booming Recently"
	thread.Description = "The Programming Language get a massive user recently."
	thread.Category.ID = "c-abc"

	if err := repo.Update(context.Background(), threadID, thread); err != nil {
		t.Fatalf("Error happened: %+v", err)
	} else {
		t.Logf("Thread with ID %s successfully updated", threadID)
	}
}

func TestDelete(t *testing.T) {
	var repo ThreadRepository = NewThreadRepositoryImpl(db)

	threadID := "t-abcdefg"

	if err := repo.Delete(context.Background(), threadID); err != nil {
		t.Fatalf("Error happened: %+v", err)
	} else {
		t.Logf("Thread with ID %s successfully deleted", threadID)
	}
}

func TestFindAllModeratorByThreadID(t *testing.T) {
	var repo ThreadRepository = NewThreadRepositoryImpl(db)

	if moderators, err := repo.FindAllModeratorByThreadID(context.Background(), "t-abcdefg"); err != nil {
		t.Fatalf("Error happened: %+v", err)
	} else {
		t.Logf("Moderators: %+v", moderators)
	}
}

func TestFindAllCommentByThreadID(t *testing.T) {
	var repo ThreadRepository = NewThreadRepositoryImpl(db)

	threadID := "t-abcdefg"
	pageInfo := entity.PageInfo{
		Limit: 20,
		Page:  1,
	}

	if paginations, err := repo.FindAllCommentByThreadID(context.Background(), threadID, pageInfo); err != nil {
		t.Fatalf("Error happened: %+v", err)
	} else {
		t.Logf("Paginations: %+v", paginations)
	}
}

func TestInsertFollowThread(t *testing.T) {
	var repo ThreadRepository = NewThreadRepositoryImpl(db)

	threadFollow := entity.ThreadFollow{
		ID: "f-1234567",
		User: entity.User{
			ID: "u-ZrxmQS",
		},
		Thread: entity.Thread{
			ID: "t-abcdefg",
		},
	}

	if err := repo.InsertFollowThread(context.Background(), threadFollow); err != nil {
		t.Fatalf("Error happened: %+v", err)
	} else {
		t.Logf("Succssfully follow thread with ID %s", threadFollow.Thread.ID)
	}
}
