package thread

import (
	"context"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/entity"
)

type ThreadRepository interface {
	Insert(ctx context.Context, thread entity.Thread) (err error)

	FindAllWithQueryAndPagination(
		ctx context.Context,
		accessorUserID string,
		query string,
		pageInfo entity.PageInfo,
	) (pagination entity.Pagination[entity.Thread], err error)

	FindAllByCategoryIDWithPagination(
		ctx context.Context,
		accessorUserID string,
		categoryID string,
		pageInfo entity.PageInfo,
	) (pagination entity.Pagination[entity.Thread], err error)

	FindAllByUserIDWithPagination(
		ctx context.Context,
		accessorUserID string,
		UserID string,
		pageInfo entity.PageInfo,
	) (pagination entity.Pagination[entity.Thread], err error)

	FindByID(
		ctx context.Context,
		accessorUserID string,
		ID string,
	) (thread entity.Thread, err error)

	Update(
		ctx context.Context,
		ID string,
		thread entity.Thread,
	) (err error)

	Delete(
		ctx context.Context,
		ID string,
	) (err error)

	FindAllModeratorByThreadID(
		ctx context.Context,
		threadID string,
	) (moderators []entity.Moderator, err error)

	FindAllCommentByThreadID(
		ctx context.Context,
		threadID string,
		pageInfo entity.PageInfo,
	) (pagination entity.Pagination[entity.Comment], err error)

	InsertComment(
		ctx context.Context,
		comment entity.Comment,
	) (err error)

	InsertFollowThread(
		ctx context.Context,
		threadFollow entity.ThreadFollow,
	) (err error)

	DeleteFollowThread(
		ctx context.Context,
		threadFollow entity.ThreadFollow,
	) (err error)

	InsertLike(
		ctx context.Context,
		like entity.Like,
	) (err error)

	DeleteLike(
		ctx context.Context,
		like entity.Like,
	) (err error)

	InsertModerator(
		ctx context.Context,
		moderator entity.Moderator,
	) (err error)

	DeleteModerator(
		ctx context.Context,
		moderator entity.Moderator,
	) (err error)

	IncrementTotalViewer(
		ID string,
	) (err error)
}
