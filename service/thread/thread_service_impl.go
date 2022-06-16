package thread

import (
	"context"
	"time"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/entity"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/response"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository/thread"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator"
)

type threadServiceImpl struct {
	threadRepository thread.ThreadRepository
	idGenerator      generator.IDGenerator
}

func NewThreadServiceImpl(
	threadRepository thread.ThreadRepository,
	idGenerator generator.IDGenerator,
) *threadServiceImpl {
	return &threadServiceImpl{
		threadRepository: threadRepository,
		idGenerator:      idGenerator,
	}
}

func (t *threadServiceImpl) GetAll(
	ctx context.Context,
	tp generator.TokenPayload,
	page uint,
	limit uint,
	query string,
) (rs response.Pagination[response.ManyThread], err error) {
	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 10
	}

	var pagination entity.Pagination[entity.Thread]
	var repoErr error

	if query == "" {
		pagination, repoErr = t.threadRepository.FindAllWithPagination(ctx, tp.ID, entity.PageInfo{Page: page, Limit: limit})
	} else {
		pagination, repoErr = t.threadRepository.FindAllWithQueryAndPagination(ctx, tp.ID, query, entity.PageInfo{Page: page, Limit: limit})
	}

	if repoErr != nil {
		err = service.MapError(repoErr)
		return
	}

	rs.PageInfo.Limit = pagination.PageInfo.Limit
	rs.PageInfo.Page = pagination.PageInfo.Page
	rs.PageInfo.PageTotal = pagination.PageInfo.PageTotal
	rs.PageInfo.Total = pagination.PageInfo.Total

	rs.List = make([]response.ManyThread, len(pagination.List))

	for i, item := range pagination.List {
		thread := response.ManyThread{
			ID:              item.ID,
			Title:           item.Title,
			CategoryID:      item.Category.ID,
			CategoryName:    item.Category.Name,
			PublishedOn:     item.CreatedAt.Format(time.RFC822),
			IsLiked:         item.IsLiked,
			IsFollowed:      item.IsFollowed,
			Description:     item.Description,
			TotalViewer:     item.TotalViewer,
			TotalLike:       item.TotalLike,
			TotalFollower:   item.TotalFollower,
			TotalComment:    item.TotalComment,
			CreatorID:       item.Creator.ID,
			CreatorUsername: item.Creator.Username,
			CreatorName:     item.Creator.Name,
		}
		rs.List[i] = thread
	}

	return
}
