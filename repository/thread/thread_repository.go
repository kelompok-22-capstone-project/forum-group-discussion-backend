package thread

import (
	"context"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/entity"
)

type ThreadRepository interface {
	Insert(ctx context.Context, thread entity.Thread) (err error)

	FindAllWithPagination(
		ctx context.Context,
		accessorUserID string,
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

	InsertFollowThread(
		ctx context.Context,
		threadFollow entity.ThreadFollow,
	) (err error)
}
