package category

import (
	"context"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/payload"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/response"
)

type CategoryService interface {
	GetAll(ctx context.Context) (rs []response.Category, err error)
	Create(ctx context.Context, accessorRole string, p payload.CreateCategory) (id string, err error)
	Update(ctx context.Context, accessorRole string, id string, p payload.UpdateCategory) (err error)
	Delete(ctx context.Context, accessorRole string, id string) (err error)
	GetAllByCategory(
		ctx context.Context,
		accessorID string,
		categoryID string,
		page uint,
		limit uint,
	) (rs response.Pagination[response.ManyThread], err error)
}
