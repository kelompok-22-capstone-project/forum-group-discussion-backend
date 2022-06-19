package user

import (
	"context"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/payload"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/response"
)

type UserService interface {
	Register(ctx context.Context, p payload.Register) (id string, err error)

	Login(ctx context.Context, p payload.Login) (r response.Login, err error)

	GetAll(
		ctx context.Context,
		accessorUserID,
		orderBy,
		status string,
		page,
		limit uint,
	) (r response.Pagination[response.User], err error)

	GetOwn(
		ctx context.Context,
		accessorUserID,
		accessorUsername string,
	) (r response.User, err error)
}
