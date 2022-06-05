package category

import (
	"context"
	"database/sql"
	"log"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/entity"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository"
)

type categoryRepositoryImpl struct {
	db *sql.DB
}

func NewCategoryRepositoryImpl(db *sql.DB) *categoryRepositoryImpl {
	return &categoryRepositoryImpl{db: db}
}

func (c *categoryRepositoryImpl) FindAll(ctx context.Context) (categories []entity.Category, err error) {
	statement := "SELECT id, name, description, created_at, updated_at FROM categories;"

	rows, dbErr := c.db.QueryContext(ctx, statement)
	if dbErr != nil {
		log.Println(dbErr)
		err = repository.ErrDatabase
		return
	}

	defer func(rows *sql.Rows) {
		if dbErr := rows.Close(); dbErr != nil {
			log.Println(dbErr)
		}
	}(rows)

	categories = make([]entity.Category, 0)
	for rows.Next() {
		var category entity.Category
		if dbErr := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.CreatedAt,
			&category.UpdatedAt,
		); dbErr != nil {
			log.Println(dbErr)
			err = repository.ErrDatabase
			return
		}
		categories = append(categories, category)
	}

	if dbErr := rows.Err(); dbErr != nil {
		log.Println(dbErr)
		err = repository.ErrDatabase
		return
	}

	return
}
func (c *categoryRepositoryImpl) Insert(ctx context.Context, category entity.Category) (err error) {
	statement := "INSERT INTO categories(id, name, description) VALUES ($1, $2, $3);"

	result, dbErr := c.db.ExecContext(ctx, statement, category.ID, category.Name, category.Description)
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
func (c *categoryRepositoryImpl) Update(ctx context.Context, ID string, category entity.Category) (err error) {
	statement := "UPDATE categories SET name = $2, description = $3, updated_at = current_timestamp WHERE id = $1;"

	result, dbErr := c.db.ExecContext(ctx, statement, ID, category.Name, category.Description)
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
		err = repository.ErrRecordNotFound
		return
	}

	return
}
func (c *categoryRepositoryImpl) Delete(ctx context.Context, ID string) (err error) {
	statement := "DELETE FROM categories WHERE id = $1;"

	result, dbErr := c.db.ExecContext(ctx, statement, ID)
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
		err = repository.ErrRecordNotFound
		return
	}

	return
}
