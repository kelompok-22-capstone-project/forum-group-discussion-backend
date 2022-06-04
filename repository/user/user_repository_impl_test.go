package user

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
	godotenv.Load("../../.env.example")
	var err error
	db, err = config.NewPostgreSQLDatabase()
	if err != nil {
		panic(err)
	}
}

func TestInsert(t *testing.T) {
	var repo UserRepository = NewUserRepositoryImpl(db)

	user := entity.User{
		ID:       "u-E8ma2R",
		Username: "erikrios",
		Email:    "erikriosetiawan15@gmail.com",
		Name:     "Erik Rio Setiawan",
		Password: "erikriosetiawan",
		Role:     "user",
	}

	if err := repo.Insert(context.Background(), user); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Successfully inserted user into database with ID %s", user.ID)
	}
}

func TestFindByUsername(t *testing.T) {
	var repo UserRepository = NewUserRepositoryImpl(db)

	if user, err := repo.FindByUsername(context.Background(), "erikrios"); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("Successfully get user: %+v", user)
	}
}
