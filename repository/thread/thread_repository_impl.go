package thread

import (
	"context"
	"database/sql"
	"log"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/entity"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository"
)

type threadRepositoryImpl struct {
	db *sql.DB
}

func NewThreadRepositoryImpl(db *sql.DB) *threadRepositoryImpl {
	return &threadRepositoryImpl{db: db}
}

func (t *threadRepositoryImpl) Insert(ctx context.Context, thread entity.Thread) (err error) {
	statement := "INSERT INTO threads(id, title, description, creator_id, category_id) VALUES ($1, $2, $3, $4, $5);"

	result, dbErr := t.db.ExecContext(ctx, statement, thread.ID, thread.Title, thread.Description, thread.Creator.ID, thread.Category.ID)
	if dbErr != nil {
		log.Println(dbErr)
		err = repository.ErrDatabase
		return
	}

	count, dbErr := result.RowsAffected()
	if dbErr != nil {
		log.Println(dbErr)
		err = repository.ErrDatabase
		return
	}

	if count < 1 {
		log.Println(dbErr)
		err = repository.ErrDatabase
		return
	}

	return
}
