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
		keyword string,
	) (pagination entity.Pagination[entity.User], err error)

	FindByUsernameWithAccessor(
		ctx context.Context,
		accessorUserID string,
		username string,
	) (user entity.User, err error)

	BannedUser(
		ctx context.Context,
		userID string,
	) (err error)

	UnbannedUser(
		ctx context.Context,
		userID string,
	) (err error)

	FollowUser(
		ctx context.Context,
		ID string,
		accessorUserID string,
		userID string,
	) (err error)

	UnfollowUser(
		ctx context.Context,
		accessorUserID string,
		userID string,
	) (err error)
}
