package controller

import (
	"net/http"
	"strconv"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service/thread"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service/user"
	"github.com/labstack/echo/v4"
)

type guestController struct {
	threadService thread.ThreadService
	userService   user.UserService
}

func NewGuestController(
	threadService thread.ThreadService,
	userService user.UserService,
) *guestController {
	return &guestController{threadService: threadService, userService: userService}
}

func (gc *guestController) Route(g *echo.Group) {
	group := g.Group("/guest")
	group.GET("/threads", gc.getThreads)
	group.GET("/threads/:id", gc.getThread)
	group.GET("/threads/:id/comments", gc.getThreadComments)
	group.GET("/users", gc.getUsers)
	group.GET("/users/:username", gc.getUserByUsername)
	group.GET("/users/:username/threads", gc.getUserThreads)
}

// getThreads     godoc
// @Summary      Get Threads
// @Description  This endpoint is used to get all threads
// @Tags         guest
// @Produce      json
// @Param        page    query  int     false  "page, default 1"
// @Param        limit   query  int     false  "limit, default 10"
// @Param        search  query  string  false  "search by keyword, default empty string"
// @Security     ApiKey
// @Success      200  {object}  threadsResponse
// @Failure      500  {object}  echo.HTTPError
// @Router       /guest/threads [get]
func (g *guestController) getThreads(c echo.Context) error {
	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")
	search := c.QueryParam("search")

	page, convErr := strconv.Atoi(pageStr)
	if convErr != nil {
		page = 0
	}

	limit, convErr := strconv.Atoi(limitStr)
	if convErr != nil {
		limit = 0
	}

	threadsResponse, err := g.threadService.GetAll(c.Request().Context(), "", uint(page), uint(limit), search)
	if err != nil {
		return newErrorResponse(err)
	}

	response := model.NewResponse("success", "Get threads successful.", threadsResponse)
	return c.JSON(http.StatusOK, response)
}

// getThread godoc
// @Summary      Get Thread by ID
// @Description  This endpoint is used to get thread by ID
// @Tags         guest
// @Produce      json
// @Param        id  path  string  true  "thread ID"
// @Security     ApiKey
// @Success      200  {object}  threadResponse
// @Failure      401  {object}  echo.HTTPError
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /guest/threads/{id} [get]
func (g *guestController) getThread(c echo.Context) error {
	id := c.Param("id")

	threadResponse, err := g.threadService.GetByID(c.Request().Context(), "", id)
	if err != nil {
		return newErrorResponse(err)
	}

	response := model.NewResponse("success", "Get thread successful.", threadResponse)
	return c.JSON(http.StatusOK, response)
}

// getThreadComments godoc
// @Summary      Get Thread Comments
// @Description  This endpoint is used to get the thread comments
// @Tags         guest
// @Produce      json
// @Param        id     path   string  true   "thread ID"
// @Param        page   query  int     false  "page, default 1"
// @Param        limit  query  int     false  "limit, default 20"
// @Security     ApiKey
// @Success      200  {object}  commentsResponse
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /guest/threads/{id}/comments [get]
func (g *guestController) getThreadComments(c echo.Context) error {
	id := c.Param("id")
	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")

	page, convErr := strconv.Atoi(pageStr)
	if convErr != nil {
		page = 0
	}

	limit, convErr := strconv.Atoi(limitStr)
	if convErr != nil {
		limit = 0
	}

	commentsResponse, err := g.threadService.GetComments(c.Request().Context(), id, uint(page), uint(limit))
	if err != nil {
		return newErrorResponse(err)
	}

	response := model.NewResponse("success", "Get comments successful", commentsResponse)
	return c.JSON(http.StatusOK, response)
}

// getUsers     godoc
// @Summary      Get Users
// @Description  This endpoint is used to get all users
// @Tags         guest
// @Produce      json
// @Param        page      query  int     false  "page, default 1"
// @Param        limit     query  int     false  "limit, default 20"
// @Param        order_by  query  string  false  "options: registered_date, ranking, default registered_date"
// @Param        status    query  string  false  "options: active, banned, default active"
// @Param        keyword   query  string  false  "search by keyword, default empty string"
// @Security     ApiKey
// @Success      200  {object}  profilesResponse
// @Failure      500  {object}  echo.HTTPError
// @Router       /guest/users [get]
func (g *guestController) getUsers(c echo.Context) error {
	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")
	orderBy := c.QueryParam("order_by")
	status := c.QueryParam("status")
	keyword := c.QueryParam("keyword")

	page, convErr := strconv.Atoi(pageStr)
	if convErr != nil {
		page = 0
	}

	limit, convErr := strconv.Atoi(limitStr)
	if convErr != nil {
		limit = 0
	}

	usersResponse, err := g.userService.GetAll(
		c.Request().Context(),
		"",
		orderBy,
		status,
		uint(page),
		uint(limit),
		keyword,
	)

	if err != nil {
		return newErrorResponse(err)
	}

	response := model.NewResponse("success", "Get users successful.", usersResponse)

	return c.JSON(http.StatusOK, response)
}

// getUserByUsername godoc
// @Summary      Get User by Username
// @Description  This endpoint is used to get the another user by username
// @Tags         guest
// @Produce      json
// @Param        username  path  string  true  "username"
// @Security     ApiKey
// @Success      200  {object}  profileResponse
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /guest/users/{username} [get]
func (g *guestController) getUserByUsername(c echo.Context) error {
	username := c.Param("username")

	userResponse, err := g.userService.GetByUsername(c.Request().Context(), "", username)
	if err != nil {
		return newErrorResponse(err)
	}

	response := model.NewResponse("success", "Get user successful.", userResponse)

	return c.JSON(http.StatusOK, response)
}

// getUserThreads godoc
// @Summary      Get User Threads
// @Description  This endpoint is used to get the user threads
// @Tags         guest
// @Produce      json
// @Param        username  path   string  true   "username"
// @Param        page      query  int     false  "page, default 1"
// @Param        limit     query  int     false  "limit, default 10"
// @Security     ApiKey
// @Success      200  {object}  threadsResponse
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /guest/users/{username}/threads [get]
func (g *guestController) getUserThreads(c echo.Context) error {
	username := c.Param("username")

	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")

	page, convErr := strconv.Atoi(pageStr)
	if convErr != nil {
		page = 0
	}

	limit, convErr := strconv.Atoi(limitStr)
	if convErr != nil {
		limit = 0
	}

	threadsResponse, err := g.userService.GetAllThreadByUsername(
		c.Request().Context(),
		"",
		username,
		uint(page),
		uint(limit),
	)

	if err != nil {
		return newErrorResponse(err)
	}

	response := model.NewResponse("success", "Get threads by username successful.", threadsResponse)

	return c.JSON(http.StatusOK, response)
}
