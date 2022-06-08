package category

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/config"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/payload"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository/category"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator"
)

var db *sql.DB
var repo category.CategoryRepository
var idGenerator generator.IDGenerator

func init() {
	godotenv.Load("./../../.env.example")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var err error
	if db, err = config.NewPostgreSQLDatabase(); err != nil {
		panic(err)
	}

	repo = category.NewCategoryRepositoryImpl(db)
	idGenerator = generator.NewNanoidIDGenerator()
}

func TestGetAll(t *testing.T) {
	var service CategoryService = NewCategoryServiceImpl(repo, idGenerator)

	if categories, err := service.GetAll(context.Background()); err != nil {
		t.Fatalf("Error happened: %s", err)
	} else {
		t.Logf("Categories: %+v", categories)
	}
}

func TestCreate(t *testing.T) {
	var service CategoryService = NewCategoryServiceImpl(repo, idGenerator)

	tp := generator.TokenPayload{
		ID:       "u-Er5Spz",
		Username: "admin",
		Role:     "admin",
		IsActive: true,
	}

	p := payload.CreateCategory{
		Name:        "Technology",
		Description: "This is a technology category description.",
	}

	if id, err := service.Create(context.Background(), tp, p); err != nil {
		t.Fatalf("Error happened: %s", err)
	} else {
		t.Logf("Successfully created a category with id %s", id)
	}
}

func TestUpdate(t *testing.T) {
	var service CategoryService = NewCategoryServiceImpl(repo, idGenerator)

	id := "c-T4m"

	tp := generator.TokenPayload{
		ID:       "u-Er5Spz",
		Username: "admin",
		Role:     "admin",
		IsActive: true,
	}

	p := payload.UpdateCategory{
		Name:        "Tech",
		Description: "This is a tech category description.",
	}

	if err := service.Update(context.Background(), tp, id, p); err != nil {
		t.Fatalf("Error happened: %s", err)
	} else {
		t.Logf("Successfully updated a category with id %s", id)
	}
}

func TestDelete(t *testing.T) {
	var service CategoryService = NewCategoryServiceImpl(repo, idGenerator)

	id := "c-T4m"

	tp := generator.TokenPayload{
		ID:       "u-Er5Spz",
		Username: "admin",
		Role:     "admin",
		IsActive: true,
	}

	if err := service.Delete(context.Background(), tp, id); err != nil {
		t.Fatalf("Error happened: %s", err)
	} else {
		t.Logf("Successfully deleted a category with id %s", id)
	}
}
