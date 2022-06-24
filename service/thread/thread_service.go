package thread

import (
	"context"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/payload"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/response"
)

type ThreadService interface {
	GetAll(
		ctx context.Context,
		accessorUserID string,
		page uint,
		limit uint,
		query string,
	) (rs response.Pagination[response.ManyThread], err error)

	Create(
		ctx context.Context,
		accessorUserID string,
		p payload.CreateThread,
	) (id string, err error)

	GetByID(
		ctx context.Context,
		accessorUserID string,
		ID string,
	) (rs response.Thread, err error)

	Update(
		ctx context.Context,
		accessorUserID string,
		ID string,
		p payload.UpdateThread,
	) (err error)

	Delete(
		ctx context.Context,
		accessorUserID string,
		role string,
		ID string,
	) (err error)

	GetComments(
		ctx context.Context,
		threadID string,
		page uint,
		limit uint,
	) (rs response.Pagination[response.Comment], err error)

	CreateComment(
		ctx context.Context,
		threadID string,
		accessorUserID string,
		p payload.CreateComment,
	) (id string, err error)

	ChangeFollowingState(
		ctx context.Context,
		threadID string,
		accessorUserID string,
	) (err error)

	ChangeLikeState(
		ctx context.Context,
		threadID string,
		accessorUserID string,
	) (err error)

	AddModerator(
		ctx context.Context,
		p payload.AddRemoveModerator,
		threadID string,
		accessorUserID string,
	) (err error)

	RemoveModerator(
		ctx context.Context,
		p payload.AddRemoveModerator,
		threadID string,
		accessorUserID string,
	) (err error)
}
