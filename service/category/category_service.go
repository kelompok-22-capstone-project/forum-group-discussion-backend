package category

import (
	"context"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/payload"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/response"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator"
)

type CategoryService interface {
	GetAll(ctx context.Context) (rs []response.Category, err error)
	Create(ctx context.Context, tp generator.TokenPayload, p payload.CreateCategory) (id string, err error)
	Update(ctx context.Context, tp generator.TokenPayload, id string, p payload.UpdateCategory) (err error)
	Delete(ctx context.Context, tp generator.TokenPayload, id string) (err error)
	GetAllByCategory(
		ctx context.Context,
		tp generator.TokenPayload,
		categoryID string,
		page uint,
		limit uint,
	) (rs response.Pagination[response.ManyThread], err error)
}
