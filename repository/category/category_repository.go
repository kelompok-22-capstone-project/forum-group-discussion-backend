package category

import (
	"context"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/entity"
)

type CategoryRepository interface {
	FindAll(ctx context.Context) (categories []entity.Category, err error)
	FindByID(ctx context.Context, ID string) (category entity.Category, err error)
	Insert(ctx context.Context, category entity.Category) (err error)
	Update(ctx context.Context, ID string, category entity.Category) (err error)
	Delete(ctx context.Context, ID string) (err error)
}
