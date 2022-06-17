package thread

import (
	"context"
	"time"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/entity"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/payload"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/response"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository/category"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository/thread"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator"
	"gopkg.in/validator.v2"
)

type threadServiceImpl struct {
	threadRepository   thread.ThreadRepository
	categoryRepository category.CategoryRepository
	idGenerator        generator.IDGenerator
}

func NewThreadServiceImpl(
	threadRepository thread.ThreadRepository,
	categoryRepository category.CategoryRepository,
	idGenerator generator.IDGenerator,
) *threadServiceImpl {
	return &threadServiceImpl{
		threadRepository:   threadRepository,
		categoryRepository: categoryRepository,
		idGenerator:        idGenerator,
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

func (t *threadServiceImpl) Create(
	ctx context.Context,
	tp generator.TokenPayload,
	p payload.CreateThread,
) (id string, err error) {
	if validateErr := validator.Validate(p); validateErr != nil {
		err = service.ErrInvalidPayload
		return
	}

	if _, repoErr := t.categoryRepository.FindByID(ctx, p.CategoryID); repoErr != nil {
		err = service.MapError(repoErr)
		return
	}

	id, genErr := t.idGenerator.GenerateThreadID()
	if genErr != nil {
		err = service.MapError(genErr)
		return
	}

	thread := entity.Thread{
		ID:          id,
		Title:       p.Title,
		Description: p.Description,
		Creator: entity.User{
			ID: tp.ID,
		},
		Category: entity.Category{
			ID: p.CategoryID,
		},
	}

	if repoErr := t.threadRepository.Insert(ctx, thread); repoErr != nil {
		err = service.MapError(repoErr)
		return
	}

	modID, genErr := t.idGenerator.GenerateModeratorID()
	if genErr != nil {
		err = service.MapError(genErr)
		return
	}

	moderator := entity.Moderator{
		ID: modID,
		User: entity.User{
			ID: tp.ID,
		},
		ThreadID: id,
	}

	if repoErr := t.threadRepository.InsertModerator(ctx, moderator); repoErr != nil {
		err = service.MapError(repoErr)
		return
	}

	return
}

func (t *threadServiceImpl) GetByID(
	ctx context.Context,
	tp generator.TokenPayload,
	ID string,
) (rs response.Thread, err error) {
	thread, repoErr := t.threadRepository.FindByID(ctx, tp.ID, ID)
	if repoErr != nil {
		err = service.MapError(repoErr)
		return
	}

	rs = response.Thread{
		ID:              thread.ID,
		Title:           thread.Title,
		CategoryID:      thread.Category.ID,
		CategoryName:    thread.Category.Name,
		PublishedOn:     thread.CreatedAt.Format(time.RFC822),
		IsLiked:         thread.IsLiked,
		IsFollowed:      thread.IsFollowed,
		Description:     thread.Description,
		TotalViewer:     thread.TotalViewer,
		TotalLike:       thread.TotalLike,
		TotalFollower:   thread.TotalFollower,
		TotalComment:    thread.TotalComment,
		CreatorID:       thread.Creator.ID,
		CreatorUsername: thread.Creator.Username,
		CreatorName:     thread.Creator.Name,
	}

	moderators, repoErr := t.threadRepository.FindAllModeratorByThreadID(ctx, ID)
	if repoErr != nil {
		err = service.MapError(repoErr)
		return
	}

	rs.Moderators = make([]response.Moderator, len(moderators))

	for i, item := range moderators {
		moderator := response.Moderator{
			ID:           item.ID,
			UserID:       item.User.ID,
			Username:     item.User.Username,
			Email:        item.User.Email,
			Name:         item.User.Name,
			Role:         item.User.Role,
			IsActive:     item.User.IsActive,
			RegisteredOn: item.CreatedAt.Format(time.RFC822),
		}
		rs.Moderators[i] = moderator
	}

	return
}

func (t *threadServiceImpl) Update(
	ctx context.Context,
	tp generator.TokenPayload,
	ID string,
	p payload.UpdateThread,
) (err error) {
	if validateErr := validator.Validate(p); validateErr != nil {
		err = service.ErrInvalidPayload
		return
	}

	if _, repoErr := t.categoryRepository.FindByID(ctx, p.CategoryID); repoErr != nil {
		err = service.MapError(repoErr)
		return
	}

	thread, repoErr := t.threadRepository.FindByID(ctx, tp.ID, ID)
	if repoErr != nil {
		err = service.MapError(repoErr)
		return
	}

	if tp.ID != thread.Creator.ID {
		err = service.ErrAccessForbidden
		return
	}

	thread.Title = p.Title
	thread.Description = p.Description
	thread.Category.ID = p.CategoryID

	if repoErr := t.threadRepository.Update(ctx, ID, thread); repoErr != nil {
		err = service.MapError(repoErr)
		return
	}

	return
}

func (t *threadServiceImpl) Delete(
	ctx context.Context,
	tp generator.TokenPayload,
	ID string,
) (err error) {
	thread, repoErr := t.threadRepository.FindByID(ctx, tp.ID, ID)
	if repoErr != nil {
		err = service.MapError(repoErr)
		return
	}

	if tp.Role != "admin" && tp.ID != thread.Creator.ID {
		err = service.ErrAccessForbidden
		return
	}

	if repoErr := t.threadRepository.Delete(ctx, ID); repoErr != nil {
		err = service.MapError(repoErr)
		return
	}

	return
}

func (t *threadServiceImpl) GetComments(
	ctx context.Context,
	threadID string,
	page uint,
	limit uint,
) (rs response.Pagination[response.Comment], err error) {
	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 20
	}

	if _, repoErr := t.threadRepository.FindByID(ctx, "", threadID); repoErr != nil {
		err = service.MapError(repoErr)
		return
	}

	pageInfo := entity.PageInfo{
		Limit: limit,
		Page:  page,
	}

	pagination, repoErr := t.threadRepository.FindAllCommentByThreadID(ctx, threadID, pageInfo)
	if repoErr != nil {
		err = service.MapError(repoErr)
		return
	}

	rs.PageInfo.Page = pagination.PageInfo.Page
	rs.PageInfo.Limit = pagination.PageInfo.Limit
	rs.PageInfo.PageTotal = pagination.PageInfo.PageTotal
	rs.PageInfo.Total = pagination.PageInfo.Total
	rs.List = make([]response.Comment, len(pagination.List))

	for i, item := range pagination.List {
		comment := response.Comment{
			ID:          item.ID,
			UserID:      item.User.ID,
			Username:    item.User.Username,
			Name:        item.User.Name,
			Comment:     item.Comment,
			PublishedOn: item.CreatedAt.Format(time.RFC822),
		}
		rs.List[i] = comment
	}

	return
}
