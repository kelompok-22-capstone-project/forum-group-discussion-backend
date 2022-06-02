package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type threadsController struct{}

func NewThreadsController() *threadsController {
	return &threadsController{}
}

func (t *threadsController) Route(g *echo.Group) {
	group := g.Group("/threads")
	group.GET("", t.getThreads)
	group.POST("", t.postCreateThread)
	group.PUT("/:id", t.putUpdateThread)
	group.DELETE("/:id", t.deleteThread)
}

// getThreads     godoc
// @Summary      Get Threads
// @Description  This endpoint is used to get all threads
// @Tags         threads
// @Produce      json
// @Success      200  {object}  threadsResponse
// @Failure      500  {object}  echo.HTTPError
// @Router       /threads [get]
func (t *threadsController) getThreads(c echo.Context) error {
	return nil
}

// postCreateThread godoc
// @Summary      Create a Thread
// @Description  This endpoint is used to create a thread
// @Tags         threads
// @Accept       json
// @Produce      json
// @Param        default  body  payload.CreateThread  true  "request body"
// @Security     ApiKeyAuth
// @Success      201  {object}  createThreadResponse
// @Failure      400  {object}  echo.HTTPError
// @Failure      401  {object}  echo.HTTPError
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /threads [post]
func (t *threadsController) postCreateThread(c echo.Context) error {
	return nil
}

// putUpdateThread godoc
// @Summary      Update a Thread
// @Description  This endpoint is used to update a thread
// @Tags         threads
// @Accept       json
// @Produce      json
// @Param        default  body  payload.UpdateThread  true  "request body"
// @Param        id       path  string                true  "thread ID"
// @Security     ApiKeyAuth
// @Success      204
// @Failure      400  {object}  echo.HTTPError
// @Failure      401  {object}  echo.HTTPError
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /threads/{id} [put]
func (t *threadsController) putUpdateThread(c echo.Context) error {
	_ = c.Param("id")
	return c.NoContent(http.StatusNoContent)
}

// deleteThread godoc
// @Summary      Delete Thread by ID
// @Description  This endpoint is used to delete a thread by ID
// @Tags         threads
// @Produce      json
// @Param        id  path  string  true  "thread ID"
// @Security     ApiKeyAuth
// @Success      204
// @Failure      404  {object}  echo.HTTPError
// @Failure      401  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /threads/{id} [delete]
func (t *threadsController) deleteThread(c echo.Context) error {
	_ = c.Param("id")
	return c.NoContent(http.StatusNoContent)
}

// createThreadResponse struct is used for swaggo to generate the API documentation, as it doesn't support generic yet.
type createThreadResponse struct {
	Status  string `json:"status" extensions:"x-order=0"`
	Message string `json:"message" extensions:"x-order=1"`
	Data    idData `json:"data" extensions:"x-order=2"`
}

type idData struct {
	ID string `json:"ID"`
}
