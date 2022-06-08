package category

import (
	"context"
	"log"
	"time"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/entity"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/payload"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/response"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/repository/category"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator"
	"gopkg.in/validator.v2"
)

type categoryServiceImpl struct {
	categoryRepository category.CategoryRepository
	idGenerator        generator.IDGenerator
}

func NewCategoryServiceImpl(
	categoryRepository category.CategoryRepository,
	idGenerator generator.IDGenerator,
) *categoryServiceImpl {
	return &categoryServiceImpl{
		categoryRepository: categoryRepository,
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
