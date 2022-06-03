package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type categoriesController struct{}

func NewCategoriesController() *categoriesController {
	return &categoriesController{}
}

func (c *categoriesController) Route(g *echo.Group) {
	group := g.Group("/categories")
	group.POST("", c.postCreateCategory)
	group.GET("", c.getCategories)
	group.PUT("/:id", c.putUpdateCategory)
	group.DELETE("/:id", c.deleteCategory)
	group.GET("/:id/threads", c.getCategoryThreads)
}

// postCreateCategory godoc
// @Summary      Create a Category
// @Description  This endpoint is used to create a category
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        default  body  payload.CreateCategory  true  "request body"
// @Security     ApiKeyAuth
// @Success      201  {object}  createThreadResponse
// @Failure      400  {object}  echo.HTTPError
// @Failure      401  {object}  echo.HTTPError
// @Failure      404  {object}  echo.HTTPError
// @Failure      500    {object}  echo.HTTPError
// @Router       /categories [post]
func (c *categoriesController) postCreateCategory(e echo.Context) error {
	return nil
}

// getCategories     godoc
// @Summary      Get Categories
// @Description  This endpoint is used to get all category
// @Tags         categories
// @Produce      json
// @Success      200  {object}  categoriesResponse
// @Failure      500  {object}  echo.HTTPError
// @Router       /categories [get]
func (c *categoriesController) getCategories(e echo.Context) error {
	return nil
}

// putUpdateCategory godoc
// @Summary      Update a Category
// @Description  This endpoint is used to update a category
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        default  body  payload.UpdateCategory  true  "request body"
// @Param        id       path  string                  true  "category ID"
// @Security     ApiKeyAuth
// @Success      204
// @Failure      400  {object}  echo.HTTPError
// @Failure      401  {object}  echo.HTTPError
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /categories/{id} [put]
func (c *categoriesController) putUpdateCategory(e echo.Context) error {
	_ = e.Param("id")
	return e.NoContent(http.StatusNoContent)
}

// deleteCategory godoc
// @Summary      Delete Category by ID
// @Description  This endpoint is used to delete a category by ID
// @Tags         categories
// @Produce      json
// @Param        id  path  string  true  "category ID"
// @Security     ApiKeyAuth
// @Success      204
// @Failure      404  {object}  echo.HTTPError
// @Failure      401  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /categories/{id} [delete]
func (c *categoriesController) deleteCategory(e echo.Context) error {
	_ = e.Param("id")
	return e.NoContent(http.StatusNoContent)
}

// getCategoryThreads godoc
// @Summary      Get Category Threads
// @Description  This endpoint is used to get the threads of particular category
// @Tags         categories
// @Produce      json
// @Param        id     path      string  true   "category ID"
// @Param        page   query     int     false  "page, default 1"
// @Param        limit  query     int     false  "limit, default 10"
// @Success      200    {object}  threadsResponse
// @Failure      404    {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /categories/{id}/threads [get]
func (c *categoriesController) getCategoryThreads(e echo.Context) error {
	_ = e.Param("id")
	_ = e.QueryParam("page")
	_ = e.QueryParam("limit")
	return nil
}

// categoriesResponse struct is used for swaggo to generate the API documentation, as it doesn't support generic yet.
type categoriesResponse struct {
	Status  string         `json:"status" extensions:"x-order=0"`
	Message string         `json:"message" extensions:"x-order=1"`
	Data    categoriesData `json:"data" extensions:"x-order=2"`
}

type categoriesData struct {
	Categories []categoryData `json:"categories" extensions:"x-order=0"`
}

type categoryData struct {
	ID          string `json:"ID" extensions:"x-order=0"`
	Name        string `json:"name" extensions:"x-order=1"`
	Description string `json:"description" extensions:"x-order=2"`
	// CreatedOn layout format: time.RFC822 (02 Jan 06 15:04 MST)
	CreatedOn string `json:"createdOn" extensions:"x-order=3"`
}
