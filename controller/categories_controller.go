package controller

import "github.com/labstack/echo/v4"

type categoriesController struct{}

func NewCategoriesController() *categoriesController {
	return &categoriesController{}
}

func (c *categoriesController) Route(g *echo.Group) {
	group := g.Group("/categories")
	group.GET("", c.getCategories)
	group.GET("/:id/threads", c.getCategoryThreads)
}

// getCategories     godoc
// @Summary      Get Categories
// @Description  This endpoint is used to get all category
// @Tags         categories
// @Produce      json
// @Success      200  {object}  categoriesResponse
// @Failure      500    {object}  echo.HTTPError
// @Router       /categories [get]
func (c *categoriesController) getCategories(e echo.Context) error {
	return nil
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

// profileResponse struct is used for swaggo to generate the API documentation, as it doesn't support generic yet.
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
