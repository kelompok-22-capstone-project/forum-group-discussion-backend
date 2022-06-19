package user

import (
	"context"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/entity"
)

type UserRepository interface {
	Insert(ctx context.Context, user entity.User) (err error)

	FindByUsername(ctx context.Context, username string) (user entity.User, err error)

	FindAllWithStatusAndPagination(
		ctx context.Context,
		accessorUserID string,
		orderBy entity.UserOrderBy,
		userStatus entity.UserStatus,
		pageInfo entity.PageInfo,
	) (pagination entity.Pagination[entity.User], err error)

	FindByUsernameWithAccessor(
		ctx context.Context,
		accessorUserID string,
		username string,
	) (user entity.User, err error)
}
