package thread

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

func TestFindAllWithQueryAndPagination(t *testing.T) {
	var repo ThreadRepository = NewThreadRepositoryImpl(db)

	pageInfo := entity.PageInfo{
		Limit: 10,
		Page:  1,
	}

	if pagination, err := repo.FindAllWithQueryAndPagination(context.Background(), "u-ZrxmQS", "Read", pageInfo); err != nil {
		t.Fatalf("Error happened: %+v", err)
	} else {
		t.Logf("Pagination: %+v", pagination)
	}
}

func TestFindAllByCategoryIDWithPagination(t *testing.T) {
	var repo ThreadRepository = NewThreadRepositoryImpl(db)

	pageInfo := entity.PageInfo{
		Limit: 10,
		Page:  1,
	}

	if pagination, err := repo.FindAllByCategoryIDWithPagination(context.Background(), "u-ZrxmQS", "c-abc", pageInfo); err != nil {
		t.Fatalf("Error happened: %+v", err)
	} else {
		t.Logf("Pagination: %+v", pagination)
	}
}

func TestFindAllByUserIDWithPagination(t *testing.T) {
	var repo ThreadRepository = NewThreadRepositoryImpl(db)

	pageInfo := entity.PageInfo{
		Limit: 10,
		Page:  1,
	}

	if pagination, err := repo.FindAllByUserIDWithPagination(context.Background(), "u-kt56R1", "u-ZrxmQS", pageInfo); err != nil {
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

func TestInsertComment(t *testing.T) {
	var repo ThreadRepository = NewThreadRepositoryImpl(db)

	comment := entity.Comment{
		ID: "c-gfedcba",
		User: entity.User{
			ID: "u-kt56R1",
		},
		Thread: entity.Thread{
			ID: "t-abcdefg",
		},
		Comment: "Good job guys, keep it simple.",
	}

	if err := repo.InsertComment(context.Background(), comment); err != nil {
		t.Fatalf("Error happened: %+v", err)
	} else {
		t.Logf("Comment with ID %s successfully inserted", comment.ID)
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

func TestDeleteFollowThread(t *testing.T) {
	var repo ThreadRepository = NewThreadRepositoryImpl(db)

	threadFollow := entity.ThreadFollow{
		User: entity.User{
			ID: "u-ZrxmQS",
		},
		Thread: entity.Thread{
			ID: "t-abcdefg",
		},
	}

	if err := repo.DeleteFollowThread(context.Background(), threadFollow); err != nil {
		t.Fatalf("Error happened: %+v", err)
	} else {
		t.Logf("Succssfully unfollow thread with ID %s", threadFollow.Thread.ID)
	}
}

func TestInsertLike(t *testing.T) {
	var repo ThreadRepository = NewThreadRepositoryImpl(db)

	like := entity.Like{}
	like.ID = "l-1234567"
	like.User.ID = "u-ZrxmQS"
	like.Thread.ID = "t-abcdefg"

	if err := repo.InsertLike(context.Background(), like); err != nil {
		t.Fatalf("Error happened: %+v", err)
	} else {
		t.Logf("Succssfully like thread with ID %s", like.Thread.ID)
	}
}

func TestDeleteLike(t *testing.T) {
	var repo ThreadRepository = NewThreadRepositoryImpl(db)

	like := entity.Like{}
	like.User.ID = "u-ZrxmQS"
	like.Thread.ID = "t-abcdefg"

	if err := repo.DeleteLike(context.Background(), like); err != nil {
		t.Fatalf("Error happened: %+v", err)
	} else {
		t.Logf("Succssfully unlike thread with ID %s", like.Thread.ID)
	}
}

func TestInsertModerator(t *testing.T) {
	var repo ThreadRepository = NewThreadRepositoryImpl(db)

	moderator := entity.Moderator{}
	moderator.ID = "m-1234"
	moderator.User.ID = "u-ZrxmQS"
	moderator.ThreadID = "t-abcdefg"

	if err := repo.InsertModerator(context.Background(), moderator); err != nil {
		t.Fatalf("Error happened: %+v", err)
	} else {
		t.Logf("Succssfully add moderator for thread with ID %s", moderator.ThreadID)
	}
}

func TestDeleteModerator(t *testing.T) {
	var repo ThreadRepository = NewThreadRepositoryImpl(db)

	moderator := entity.Moderator{}
	moderator.User.ID = "u-ZrxmQS"
	moderator.ThreadID = "t-abcdefg"

	if err := repo.DeleteModerator(context.Background(), moderator); err != nil {
		t.Fatalf("Error happened: %+v", err)
	} else {
		t.Logf("Succssfully delete moderator for thread with ID %s", moderator.ThreadID)
	}
}
