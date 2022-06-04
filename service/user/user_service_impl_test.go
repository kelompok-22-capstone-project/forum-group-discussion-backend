package user

import (
	"context"
	"database/sql"
	"testing"

	"github.com/joho/godotenv"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/config"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/payload"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository/user"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator"
)

var db *sql.DB
var repo user.UserRepository
var idGenerator generator.IDGenerator
var passwordGenerator generator.PasswordGenerator
var tokenGenerator generator.TokenGenerator

func init() {
	godotenv.Load("./../../.env.example")
	var err error
	if db, err = config.NewPostgreSQLDatabase(); err != nil {
		panic(err)
	}

	repo = user.NewUserRepositoryImpl(db)
	idGenerator = generator.NewNanoidIDGenerator()
	passwordGenerator = generator.NewBcryptPasswordGenerator()
	tokenGenerator = generator.NewJWTTokenGenerator()
}

func TestRegister(t *testing.T) {
	var service UserService = NewUserServiceImpl(repo, idGenerator, passwordGenerator, tokenGenerator)

	p := payload.Register{
		Username: "erikrios",
		Email:    "erikriosetiawan15@gmail.com",
		Name:     "Erik Rio Setiawan",
		Password: "erikriosetiawan",
	}

	if id, err := service.Register(context.Background(), p); err != nil {
		t.Fatalf("Error happened: %s", err)
	} else {
		t.Logf("Successfully registered user with id %s", id)
	}
}

func TestLogin(t *testing.T) {
	var service UserService = NewUserServiceImpl(repo, idGenerator, passwordGenerator, tokenGenerator)

	p := payload.Login{
		Username: "erikrios",
		Password: "erikriosetiawan",
	}

	if resp, err := service.Login(context.Background(), p); err != nil {
		t.Fatalf("Error happened: %s", err)
	} else {
		t.Logf("Successfully login with token %s and role %s", resp.Token, resp.Role)
	}
}
