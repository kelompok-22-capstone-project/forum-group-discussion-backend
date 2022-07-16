package controller

import (
	"net/http"
	"strconv"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/middleware"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/payload"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/response"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service/thread"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator"
	"github.com/labstack/echo/v4"
)

type threadsController struct {
	threadService  thread.ThreadService
	tokenGenerator generator.TokenGenerator
}

func NewThreadsController(
	threadService thread.ThreadService,
	tokenGenerator generator.TokenGenerator,
) *threadsController {
	return &threadsController{
		threadService:  threadService,
		tokenGenerator: tokenGenerator,
	}
}

func (t *threadsController) Route(g *echo.Group) {
	group := g.Group("/threads")
	group.GET("", t.getThreads, middleware.JWTMiddleware())
	group.POST("", t.postCreateThread, middleware.JWTMiddleware())
	group.GET("/:id", t.getThread, middleware.JWTMiddleware())
	group.PUT("/:id", t.putUpdateThread, middleware.JWTMiddleware())
	group.DELETE("/:id", t.deleteThread, middleware.JWTMiddleware())
	group.GET("/:id/comments", t.getThreadComments, middleware.JWTMiddleware())
	group.POST("/:id/comments", t.postCreateThreadComments, middleware.JWTMiddleware())
	group.PUT("/:id/like", t.putThreadLike, middleware.JWTMiddleware())
	group.PUT("/:id/follow", t.putThreadFollow, middleware.JWTMiddleware())
	group.PUT("/:id/moderators/add", t.putThreadAddModerator, middleware.JWTMiddleware())
	group.PUT("/:id/moderators/remove", t.putThreadRemoveModerator, middleware.JWTMiddleware())
}

// getThreads     godoc
// @Summary      Get Threads
// @Description  This endpoint is used to get all threads
// @Tags         threads
// @Produce      json
// @Param        page    query  int     false  "page, default 1"
// @Param        limit   query  int     false  "limit, default 10"
// @Param        search  query  string  false  "search by keyword, default empty string"
// @Security     ApiKey
// @Security     ApiKeyAuth
// @Success      200  {object}  threadsResponse
// @Failure      500  {object}  echo.HTTPError
// @Router       /threads [get]
func (t *threadsController) getThreads(c echo.Context) error {
	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")
	search := c.QueryParam("search")

	page, convErr := strconv.Atoi(pageStr)
	if convErr != nil || page < 0 {
		page = 0
	}

	limit, convErr := strconv.Atoi(limitStr)
	if convErr != nil || limit < 0 {
		limit = 0
	}

	tp := t.tokenGenerator.ExtractToken(c)

	threadsResponse, err := t.threadService.GetAll(c.Request().Context(), tp.ID, uint(page), uint(limit), search)
	if err != nil {
		return newErrorResponse(err)
	}

	response := model.NewResponse("success", "Get threads successful.", threadsResponse)
	return c.JSON(http.StatusOK, response)
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
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /threads [post]
func (t *threadsController) postCreateThread(c echo.Context) error {
	tp := t.tokenGenerator.ExtractToken(c)

	p := new(payload.CreateThread)
	if err := c.Bind(p); err != nil {
		return newErrorResponse(service.ErrInvalidPayload)
	}

	id, err := t.threadService.Create(c.Request().Context(), tp.ID, *p)
	if err != nil {
		return newErrorResponse(err)
	}

	idResponse := map[string]any{"ID": id}
	response := model.NewResponse("success", "Create thread successful.", idResponse)
	return c.JSON(http.StatusCreated, response)
}

// getThread godoc
// @Summary      Get Thread by ID
// @Description  This endpoint is used to get thread by ID
// @Tags         threads
// @Produce      json
// @Param        id       path  string                      true  "thread ID"
// @Security     ApiKey
// @Security     ApiKeyAuth
// @Success      200  {object}  threadResponse
// @Failure      401  {object}  echo.HTTPError
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /threads/{id} [get]
func (t *threadsController) getThread(c echo.Context) error {
	id := c.Param("id")

	tp := t.tokenGenerator.ExtractToken(c)

	threadResponse, err := t.threadService.GetByID(c.Request().Context(), tp.ID, id)
	if err != nil {
		return newErrorResponse(err)
	}

	response := model.NewResponse("success", "Get thread successful.", threadResponse)
	return c.JSON(http.StatusOK, response)
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
	id := c.Param("id")

	tp := t.tokenGenerator.ExtractToken(c)

	p := new(payload.UpdateThread)
	if err := c.Bind(p); err != nil {
		return newErrorResponse(service.ErrInvalidPayload)
	}

	if err := t.threadService.Update(c.Request().Context(), tp.ID, id, *p); err != nil {
		return newErrorResponse(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// deleteThread godoc
// @Summary      Delete Thread by ID
// @Description  This endpoint is used to delete a thread by ID
// @Tags         threads
// @Produce      json
// @Param        id  path  string  true  "thread ID"
// @Security     ApiKey
// @Security     ApiKeyAuth
// @Success      204
// @Failure      404  {object}  echo.HTTPError
// @Failure      401  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /threads/{id} [delete]
func (t *threadsController) deleteThread(c echo.Context) error {
	id := c.Param("id")

	tp := t.tokenGenerator.ExtractToken(c)

	if err := t.threadService.Delete(c.Request().Context(), tp.ID, tp.Role, id); err != nil {
		return newErrorResponse(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// getThreadComments godoc
// @Summary      Get Thread Comments
// @Description  This endpoint is used to get the thread comments
// @Tags         threads
// @Produce      json
// @Param        id       path  string                 true  "thread ID"
// @Param        page   query  int     false  "page, default 1"
// @Param        limit  query  int     false  "limit, default 20"
// @Security     ApiKey
// @Security     ApiKeyAuth
// @Success      200  {object}  commentsResponse
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /threads/{id}/comments [get]
func (t *threadsController) getThreadComments(c echo.Context) error {
	id := c.Param("id")
	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")

	page, convErr := strconv.Atoi(pageStr)
	if convErr != nil || page < 0 {
		page = 0
	}

	limit, convErr := strconv.Atoi(limitStr)
	if convErr != nil || limit < 0 {
		limit = 0
	}

	commentsResponse, err := t.threadService.GetComments(c.Request().Context(), id, uint(page), uint(limit))
	if err != nil {
		return newErrorResponse(err)
	}

	response := model.NewResponse("success", "Get comments successful", commentsResponse)
	return c.JSON(http.StatusOK, response)
}

// postCreateThreadComments godoc
// @Summary      Create a Comment
// @Description  This endpoint is used to create a comment of a thread
// @Tags         threads
// @Accept       json
// @Produce      json
// @Param        id     path   string  true   "thread ID"
// @Param        default  body  payload.CreateComment  true  "request body"
// @Security     ApiKey
// @Security     ApiKeyAuth
// @Success      201  {object}  createThreadResponse
// @Failure      400  {object}  echo.HTTPError
// @Failure      401  {object}  echo.HTTPError
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /threads/{id}/comments [post]
func (t *threadsController) postCreateThreadComments(c echo.Context) error {
	id := c.Param("id")

	tp := t.tokenGenerator.ExtractToken(c)

	p := new(payload.CreateComment)
	if err := c.Bind(p); err != nil {
		return newErrorResponse(service.ErrInvalidPayload)
	}

	id, err := t.threadService.CreateComment(c.Request().Context(), id, tp.ID, *p)
	if err != nil {
		return newErrorResponse(err)
	}

	idResponse := map[string]any{"ID": id}
	response := model.NewResponse("success", "Create comment successful.", idResponse)
	return c.JSON(http.StatusCreated, response)
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
	threadID := c.Param("id")

	tp := t.tokenGenerator.ExtractToken(c)

	if err := t.threadService.ChangeLikeState(c.Request().Context(), threadID, tp.ID); err != nil {
		return newErrorResponse(err)
	}

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
	threadID := c.Param("id")

	tp := t.tokenGenerator.ExtractToken(c)

	if err := t.threadService.ChangeFollowingState(c.Request().Context(), threadID, tp.ID); err != nil {
		return newErrorResponse(err)
	}

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
	threadID := c.Param("id")
	tp := t.tokenGenerator.ExtractToken(c)

	p := new(payload.AddRemoveModerator)
	if err := c.Bind(p); err != nil {
		return newErrorResponse(service.ErrInvalidPayload)
	}

	if err := t.threadService.AddModerator(c.Request().Context(), *p, threadID, tp.ID); err != nil {
		return newErrorResponse(err)
	}

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
	threadID := c.Param("id")
	tp := t.tokenGenerator.ExtractToken(c)

	p := new(payload.AddRemoveModerator)
	if err := c.Bind(p); err != nil {
		return newErrorResponse(service.ErrInvalidPayload)
	}

	if err := t.threadService.RemoveModerator(c.Request().Context(), *p, threadID, tp.ID); err != nil {
		return newErrorResponse(err)
	}

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
	Status  string          `json:"status" extensions:"x-order=0"`
	Message string          `json:"message" extensions:"x-order=1"`
	Data    response.Thread `json:"data" extensions:"x-order=2"`
}

// commentsResponse struct is used for swaggo to generate the API documentation, as it doesn't support generic yet.
type commentsResponse struct {
	Status  string              `json:"status" extensions:"x-order=0"`
	Message string              `json:"message" extensions:"x-order=1"`
	Data    commentsInfoWrapper `json:"data" extensions:"x-order=2"`
}

type commentsInfoWrapper struct {
	Threads  []response.Comment `json:"list" extensions:"x-order=0"`
	PageInfo pageInfoData       `json:"pageInfo" extensions:"x-order=1"`
}
