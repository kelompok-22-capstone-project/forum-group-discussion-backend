package category

import (
	"context"
	"database/sql"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/entity"
)

type categoryRepositoryImpl struct {
	db *sql.DB
}

func NewCategoryRepositoryImpl(db *sql.DB) *categoryRepositoryImpl {
	return &categoryRepositoryImpl{db: db}
}

func (c *categoryRepositoryImpl) FindAll(ctx context.Context) (categories []entity.Category, err error) {
	return
}
func (c *categoryRepositoryImpl) Insert(ctx context.Context, category entity.Category) (err error) {
	return
}
func (c *categoryRepositoryImpl) Update(ctx context.Context, ID string, category entity.Category) (err error) {
	return
}
func (c *categoryRepositoryImpl) Delete(ctx context.Context, ID string) (err error) { return }
