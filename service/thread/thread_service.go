package thread

import (
	"context"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/payload"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/response"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator"
)

type ThreadService interface {
	GetAll(
		ctx context.Context,
		tp generator.TokenPayload,
		page uint,
		limit uint,
		query string,
	) (rs response.Pagination[response.ManyThread], err error)

	Create(
		ctx context.Context,
		tp generator.TokenPayload,
		p payload.CreateThread,
	) (id string, err error)

	GetByID(
		ctx context.Context,
		tp generator.TokenPayload,
		ID string,
	) (rs response.Thread, err error)

	Update(
		ctx context.Context,
		tp generator.TokenPayload,
		ID string,
		p payload.UpdateThread,
	) (err error)

	Delete(
		ctx context.Context,
		tp generator.TokenPayload,
		ID string,
	) (err error)

	GetComments(
		ctx context.Context,
		threadID string,
	) (rs response.Pagination[response.Comment], err error)
}
