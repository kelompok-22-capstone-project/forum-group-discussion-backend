package category

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

func TestFindAll(t *testing.T) {
	var repo CategoryRepository = NewCategoryRepositoryImpl(db)

	if categories, err := repo.FindAll(context.Background()); err != nil {
		t.Fatalf("Error happened: %s", err)
	} else {
		t.Logf("Categories: %+v", categories)
	}
}

func TestFindByID(t *testing.T) {
	var repo CategoryRepository = NewCategoryRepositoryImpl(db)

	if category, err := repo.FindByID(context.Background(), "c-abc"); err != nil {
		t.Fatalf("Error happened: %s", err)
	} else {
		t.Logf("Categories: %+v", category)
	}
}

func TestInsert(t *testing.T) {
	var repo CategoryRepository = NewCategoryRepositoryImpl(db)

	category := entity.Category{
		ID:          "c-xyz",
		Name:        "Technology",
		Description: "This is technology category",
	}

	if err := repo.Insert(context.Background(), category); err != nil {
		t.Fatalf("Error happened: %s", err)
	} else {
		t.Logf("Successfully inserted category with id %s", category.ID)
	}
}

func TestUpdate(t *testing.T) {
	var repo CategoryRepository = NewCategoryRepositoryImpl(db)

	id := "c-xyz"
	category := entity.Category{
		Name:        "Tech",
		Description: "This is tech category",
	}

	if err := repo.Update(context.Background(), id, category); err != nil {
		t.Fatalf("Error happened: %s", err)
	} else {
		t.Logf("Successfully updated category with id %s", id)
	}
}

func TestDelete(t *testing.T) {

	var repo CategoryRepository = NewCategoryRepositoryImpl(db)

	id := "c-xyz"

	if err := repo.Delete(context.Background(), id); err != nil {
		t.Fatalf("Error happened: %s", err)
	} else {
		t.Logf("Successfully deleted category with id %s", id)
	}
}
