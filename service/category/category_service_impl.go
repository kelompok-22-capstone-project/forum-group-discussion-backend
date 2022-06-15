package category

import (
	"context"
	"log"
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

type categoryServiceImpl struct {
	categoryRepository category.CategoryRepository
	threadRepository   thread.ThreadRepository
	idGenerator        generator.IDGenerator
}

func NewCategoryServiceImpl(
	categoryRepository category.CategoryRepository,
	threadRepository thread.ThreadRepository,
	idGenerator generator.IDGenerator,
) *categoryServiceImpl {
	return &categoryServiceImpl{
		categoryRepository: categoryRepository,
		threadRepository:   threadRepository,
		idGenerator:        idGenerator,
	}
}

func (c *categoryServiceImpl) GetAll(ctx context.Context) (rs []response.Category, err error) {
	categories, repoErr := c.categoryRepository.FindAll(ctx)
	if repoErr != nil {
		err = service.MapError(repoErr)
		return
	}

	rs = make([]response.Category, len(categories))
	for i, category := range categories {
		response := response.Category{
			ID:          category.ID,
			Name:        category.Name,
			Description: category.Description,
			CreatedOn:   category.CreatedAt.Format(time.RFC822),
		}
		rs[i] = response
	}

	return
}

func (c *categoryServiceImpl) Create(ctx context.Context, tp generator.TokenPayload, p payload.CreateCategory) (id string, err error) {
	log.Println(tp)
	if tp.Role != "admin" {
		err = service.ErrAccessForbidden
		return
	}

	if validateErr := validator.Validate(p); validateErr != nil {
		err = service.ErrInvalidPayload
		return
	}

	id, genErr := c.idGenerator.GenerateCategoryID()
	if genErr != nil {
		err = service.MapError(genErr)
		return
	}

	category := entity.Category{
		ID:          id,
		Name:        p.Name,
		Description: p.Description,
	}

	if repoErr := c.categoryRepository.Insert(ctx, category); repoErr != nil {
		err = service.MapError(repoErr)
		return
	}

	return
}

func (c *categoryServiceImpl) Update(ctx context.Context, tp generator.TokenPayload, id string, p payload.UpdateCategory) (err error) {
	log.Println(tp)
	if tp.Role != "admin" {
		err = service.ErrAccessForbidden
		return
	}

	if validateErr := validator.Validate(p); validateErr != nil {
		err = service.ErrInvalidPayload
		return
	}

	category := entity.Category{
		Name:        p.Name,
		Description: p.Description,
	}

	if repoErr := c.categoryRepository.Update(ctx, id, category); repoErr != nil {
		err = service.MapError(repoErr)
		return
	}

	return
}

func (c *categoryServiceImpl) Delete(ctx context.Context, tp generator.TokenPayload, id string) (err error) {
	log.Println(tp)
	if tp.Role != "admin" {
		err = service.ErrAccessForbidden
		return
	}

	if repoErr := c.categoryRepository.Delete(ctx, id); repoErr != nil {
		err = service.MapError(repoErr)
		return
	}

	return
}

func (c *categoryServiceImpl) GetAllByCategory(
	ctx context.Context,
	tp generator.TokenPayload,
	categoryID string,
	page uint,
	limit uint,
) (rs response.Pagination[response.ManyThread], err error) {
	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 10
	}

	if _, repoErr := c.categoryRepository.FindByID(ctx, categoryID); repoErr != nil {
		err = service.MapError(repoErr)
		return
	}

	pagination, repoErr := c.threadRepository.FindAllByCategoryIDWithPagination(
		ctx,
		tp.ID,
		categoryID,
		entity.PageInfo{
			Limit: limit,
			Page:  page,
		},
	)

	if repoErr != nil {
		err = service.MapError(repoErr)
		return
	}

	rs.PageInfo.Page = pagination.PageInfo.Page
	rs.PageInfo.Limit = pagination.PageInfo.Limit
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
