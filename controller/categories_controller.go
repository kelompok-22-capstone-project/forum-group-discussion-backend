package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/middleware"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/payload"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/response"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service/category"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator"
	"github.com/labstack/echo/v4"
)

type categoriesController struct {
	categoryService category.CategoryService
	tokenGenerator  generator.TokenGenerator
}

func NewCategoriesController(
	categoryService category.CategoryService,
	tokenGenerator generator.TokenGenerator,
) *categoriesController {
	return &categoriesController{
		categoryService: categoryService,
		tokenGenerator:  tokenGenerator,
	}
}

func (c *categoriesController) Route(g *echo.Group) {
	group := g.Group("/categories")
	group.POST("", c.postCreateCategory, middleware.JWTMiddleware())
	group.GET("", c.getCategories)
	group.PUT("/:id", c.putUpdateCategory, middleware.JWTMiddleware())
	group.DELETE("/:id", c.deleteCategory, middleware.JWTMiddleware())
	group.GET("/:id/threads", c.getCategoryThreads, middleware.JWTMiddleware())
}

// postCreateCategory godoc
// @Summary      Create a Category
// @Description  This endpoint is used to create a category
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        default  body  payload.CreateCategory  true  "request body"
// @Security     ApiKey
// @Security     ApiKeyAuth
// @Success      201  {object}  createThreadResponse
// @Failure      400  {object}  echo.HTTPError
// @Failure      401  {object}  echo.HTTPError
// @Failure      403  {object}  echo.HTTPError
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /categories [post]
func (c *categoriesController) postCreateCategory(e echo.Context) error {
	tp := c.tokenGenerator.ExtractToken(e)

	p := new(payload.CreateCategory)
	if err := e.Bind(p); err != nil {
		return newErrorResponse(service.ErrInvalidPayload)
	}

	id, err := c.categoryService.Create(e.Request().Context(), tp.ID, *p)
	if err != nil {
		return newErrorResponse(err)
	}

	idResponse := map[string]any{"ID": id}
	response := model.NewResponse("success", "Create category successful.", idResponse)
	return e.JSON(http.StatusCreated, response)
}

// getCategories     godoc
// @Summary      Get Categories
// @Description  This endpoint is used to get all category
// @Tags         categories
// @Produce      json
// @Security     ApiKey
// @Success      200  {object}  categoriesResponse
// @Failure      500  {object}  echo.HTTPError
// @Router       /categories [get]
func (c *categoriesController) getCategories(e echo.Context) error {
	tp := c.tokenGenerator.ExtractToken(e)
	log.Printf("%+v", tp)
	categories, err := c.categoryService.GetAll(e.Request().Context())
	if err != nil {
		return newErrorResponse(err)
	}

	categoriesResponse := map[string]any{"categories": categories}
	response := model.NewResponse("success", "Get categories successful.", categoriesResponse)
	return e.JSON(http.StatusOK, response)
}

// putUpdateCategory godoc
// @Summary      Update a Category
// @Description  This endpoint is used to update a category
// @Tags         categories
// @Accept       json
// @Produce      json
// @Security     ApiKey
// @Param        default  body  payload.UpdateCategory  true  "request body"
// @Param        id       path  string                  true  "category ID"
// @Security     ApiKeyAuth
// @Success      204
// @Failure      400  {object}  echo.HTTPError
// @Failure      401  {object}  echo.HTTPError
// @Failure      403  {object}  echo.HTTPError
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /categories/{id} [put]
func (c *categoriesController) putUpdateCategory(e echo.Context) error {
	id := e.Param("id")

	tp := c.tokenGenerator.ExtractToken(e)

	p := new(payload.UpdateCategory)
	if err := e.Bind(p); err != nil {
		return newErrorResponse(service.ErrInvalidPayload)
	}

	if err := c.categoryService.Update(e.Request().Context(), tp.ID, id, *p); err != nil {
		return newErrorResponse(err)
	}

	return e.NoContent(http.StatusNoContent)
}

// deleteCategory godoc
// @Summary      Delete Category by ID
// @Description  This endpoint is used to delete a category by ID
// @Tags         categories
// @Produce      json
// @Param        id  path  string  true  "category ID"
// @Security     ApiKey
// @Security     ApiKeyAuth
// @Success      204
// @Failure      404  {object}  echo.HTTPError
// @Failure      401  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /categories/{id} [delete]
func (c *categoriesController) deleteCategory(e echo.Context) error {
	id := e.Param("id")

	tp := c.tokenGenerator.ExtractToken(e)

	if err := c.categoryService.Delete(e.Request().Context(), tp.ID, id); err != nil {
		return newErrorResponse(err)
	}

	return e.NoContent(http.StatusNoContent)
}

// getCategoryThreads godoc
// @Summary      Get Category Threads
// @Description  This endpoint is used to get the threads of particular category
// @Tags         categories
// @Produce      json
// @Param        id     path   string  true   "category ID"
// @Param        page   query  int     false  "page, default 1"
// @Param        limit  query  int     false  "limit, default 10"
// @Security     ApiKey
// @Security     ApiKeyAuth
// @Success      200  {object}  threadsResponse
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /categories/{id}/threads [get]
func (c *categoriesController) getCategoryThreads(e echo.Context) error {
	id := e.Param("id")
	pageStr := e.QueryParam("page")
	limitStr := e.QueryParam("limit")

	page, convErr := strconv.Atoi(pageStr)
	if convErr != nil {
		page = 0
	}

	limit, convErr := strconv.Atoi(limitStr)
	if convErr != nil {
		limit = 0
	}

	tp := c.tokenGenerator.ExtractToken(e)

	threadsResponse, err := c.categoryService.GetAllByCategory(e.Request().Context(), tp.ID, id, uint(page), uint(limit))
	if err != nil {
		return newErrorResponse(err)
	}

	response := model.NewResponse("success", "Get threads by category successful.", threadsResponse)
	return e.JSON(http.StatusOK, response)
}

// categoriesResponse struct is used for swaggo to generate the API documentation, as it doesn't support generic yet.
type categoriesResponse struct {
	Status  string         `json:"status" extensions:"x-order=0"`
	Message string         `json:"message" extensions:"x-order=1"`
	Data    categoriesData `json:"data" extensions:"x-order=2"`
}

type categoriesData struct {
	Categories []response.Category `json:"categories" extensions:"x-order=0"`
}
