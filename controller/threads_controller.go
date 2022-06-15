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
	group.GET("/:id", t.getThread)
	group.PUT("/:id", t.putUpdateThread)
	group.DELETE("/:id", t.deleteThread)
	group.GET("/:id/comments", t.getThreadComments)
	group.PUT("/:id/like", t.putThreadLike)
	group.PUT("/:id/follow", t.putThreadFollow)
	group.PUT("/:id/moderators/add", t.putThreadAddModerator)
	group.PUT("/:id/moderators/remove", t.putThreadRemoveModerator)
}

// getThreads     godoc
// @Summary      Get Threads
// @Description  This endpoint is used to get all threads
// @Tags         threads
// @Produce      json
// @Param        page    query     int     false  "page, default 1"
// @Param        limit   query     int     false  "limit, default 10"
// @Param        search  query     string  false  "search by keyword, default empty string"
// @Security     ApiKey
// @Success      200     {object}  threadsResponse
// @Failure      500     {object}  echo.HTTPError
// @Router       /threads [get]
func (t *threadsController) getThreads(c echo.Context) error {
	_ = c.QueryParam("page")
	_ = c.QueryParam("limit")
	return nil
}

// postCreateThread godoc
// @Summary      Create a Thread
// @Description  This endpoint is used to create a thread
// @Tags         threads
// @Accept       json
// @Produce      json
// @Param        default  body  payload.CreateThread  true  "request body"
// @Security     ApiKey
// @Security     ApiKeyAuth
// @Success      201  {object}  createThreadResponse
// @Failure      400  {object}  echo.HTTPError
// @Failure      401  {object}  echo.HTTPError
// @Failure      404    {object}  echo.HTTPError
// @Failure      500    {object}  echo.HTTPError
// @Router       /threads [post]
func (t *threadsController) postCreateThread(c echo.Context) error {
	return nil
}

// getThread godoc
// @Summary      Get Thread by ID
// @Description  This endpoint is used to get thread by ID
// @Tags         threads
// @Produce      json
// @Param        id   path      string  true  "thread ID"
// @Security     ApiKey
// @Success      200  {object}  threadResponse
// @Failure      401  {object}  echo.HTTPError
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /threads/{id} [get]
func (t *threadsController) getThread(c echo.Context) error {
	_ = c.Param("id")

	return nil
}

// putUpdateThread godoc
// @Summary      Update a Thread
// @Description  This endpoint is used to update a thread
// @Tags         threads
// @Accept       json
// @Produce      json
// @Param        default  body  payload.UpdateThread  true  "request body"
// @Param        id       path  string                      true  "thread ID"
// @Security     ApiKey
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
// @Param        id       path  string                      true  "thread ID"
// @Security     ApiKey
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

// getThreadComments godoc
// @Summary      Get Thread Comments
// @Description  This endpoint is used to get the thread comments
// @Tags         threads
// @Produce      json
// @Param        id     path      string  true   "thread ID"
// @Param        page   query     int     false  "page, default 1"
// @Param        limit  query     int     false  "limit, default 20"
// @Security     ApiKey
// @Success      200    {object}  commentsResponse
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /threads/{id}/comments [get]
func (t *threadsController) getThreadComments(c echo.Context) error {
	_ = c.Param("id")
	_ = c.QueryParam("page")
	_ = c.QueryParam("limit")
	return nil
}

// putThreadLike godoc
// @Summary      Like/Unlike a Thread
// @Description  This endpoint is used to like/unlike a thread
// @Tags         threads
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "thread ID"
// @Security     ApiKey
// @Security     ApiKeyAuth
// @Success      204
// @Failure      401  {object}  echo.HTTPError
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /threads/{id}/like [put]
func (t *threadsController) putThreadLike(c echo.Context) error {
	_ = c.Param("id")
	return c.NoContent(http.StatusNoContent)
}

// putThreadFollow godoc
// @Summary      Follow/Unfollow a Thread
// @Description  This endpoint is used to follow/unfollow a thread
// @Tags         threads
// @Accept       json
// @Produce      json
// @Param        id       path  string                true  "thread ID"
// @Security     ApiKey
// @Security     ApiKeyAuth
// @Success      204
// @Failure      401  {object}  echo.HTTPError
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /threads/{id}/follow [put]
func (t *threadsController) putThreadFollow(c echo.Context) error {
	_ = c.Param("id")
	return c.NoContent(http.StatusNoContent)
}

// putThreadAddModerator godoc
// @Summary      Add a Moderator to Thread
// @Description  This endpoint is used to add a moderator to thread
// @Tags         threads
// @Accept       json
// @Produce      json
// @Param        default  body  payload.AddRemoveModerator  true  "request body"
// @Param        id  path  string  true  "thread ID"
// @Security     ApiKey
// @Security     ApiKeyAuth
// @Success      204
// @Failure      400  {object}  echo.HTTPError
// @Failure      401  {object}  echo.HTTPError
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /threads/{id}/moderators/add [put]
func (t *threadsController) putThreadAddModerator(c echo.Context) error {
	_ = c.Param("id")
	return c.NoContent(http.StatusNoContent)
}

// putThreadRemoveModerator godoc
// @Summary      Remove a Moderator from Thread
// @Description  This endpoint is used to remove a moderator from thread
// @Tags         threads
// @Accept       json
// @Produce      json
// @Param        default  body  payload.AddRemoveModerator  true  "request body"
// @Param        id  path  string  true  "thread ID"
// @Security     ApiKey
// @Security     ApiKeyAuth
// @Success      204
// @Failure      400  {object}  echo.HTTPError
// @Failure      401  {object}  echo.HTTPError
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /threads/{id}/moderators/remove [put]
func (t *threadsController) putThreadRemoveModerator(c echo.Context) error {
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

// threadResponse struct is used for swaggo to generate the API documentation, as it doesn't support generic yet.
type threadResponse struct {
	Status  string     `json:"status" extensions:"x-order=0"`
	Message string     `json:"message" extensions:"x-order=1"`
	Data    threadData `json:"data" extensions:"x-order=2"`
}

// commentsResponse struct is used for swaggo to generate the API documentation, as it doesn't support generic yet.
type commentsResponse struct {
	Status  string              `json:"status" extensions:"x-order=0"`
	Message string              `json:"message" extensions:"x-order=1"`
	Data    commentsInfoWrapper `json:"data" extensions:"x-order=2"`
}

type commentsInfoWrapper struct {
	Threads  []commentData `json:"comments" extensions:"x-order=0"`
	PageInfo pageInfoData  `json:"pageInfo" extensions:"x-order=1"`
}

type commentData struct {
	ID       string `json:"ID" extensions:"x-order=0"`
	UserID   string `json:"userID" extensions:"x-order=1"`
	Username string `json:"username" extensions:"x-order=2"`
	Name     string `json:"name" extensions:"x-order=3"`
	Comment  string `json:"comment" extensions:"x-order=4"`
	// PublishedOn layout format: time.RFC822 (02 Jan 06 15:04 MST)
	PublishedOn string `json:"publishedOn" extensions:"x-order=5"`
}
