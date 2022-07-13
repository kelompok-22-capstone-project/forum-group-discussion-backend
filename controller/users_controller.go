package controller

import (
	"net/http"
	"strconv"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/middleware"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/response"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service/user"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator"
	"github.com/labstack/echo/v4"
)

type usersController struct {
	userService    user.UserService
	tokenGenerator generator.TokenGenerator
}

func NewUsersController(
	userService user.UserService,
	tokenGenerator generator.TokenGenerator,
) *usersController {
	return &usersController{
		userService:    userService,
		tokenGenerator: tokenGenerator,
	}
}

func (u *usersController) Route(g *echo.Group) {
	group := g.Group("/users")
	group.GET("", u.getUsers, middleware.JWTMiddleware())
	group.GET("/me", u.getMe, middleware.JWTMiddleware())
	group.GET("/:username", u.getUserByUsername, middleware.JWTMiddleware())
	group.GET("/:username/threads", u.getUserThreads, middleware.JWTMiddleware())
	group.PUT("/:username/follow", u.putUserFollow, middleware.JWTMiddleware())
	group.PUT("/:username/banned", u.putUserBanned, middleware.JWTMiddleware())
}

// getUsers     godoc
// @Summary      Get Users
// @Description  This endpoint is used to get all users
// @Tags         users
// @Produce      json
// @Param        page      query  int     false  "page, default 1"
// @Param        limit     query  int     false  "limit, default 20"
// @Param        order_by  query  string  false  "options: registered_date, ranking, default registered_date"
// @Param        status    query  string  false  "options: active, banned, default active"
// @Param        keyword   query  string  false  "search by keyword, default empty string"
// @Security     ApiKey
// @Security     ApiKeyAuth
// @Success      200  {object}  profilesResponse
// @Failure      500  {object}  echo.HTTPError
// @Router       /users [get]
func (u *usersController) getUsers(c echo.Context) error {
	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")
	orderBy := c.QueryParam("order_by")
	status := c.QueryParam("status")
	keyword := c.QueryParam("keyword")

	page, convErr := strconv.Atoi(pageStr)
	if convErr != nil || page < 0 {
		page = 0
	}

	limit, convErr := strconv.Atoi(limitStr)
	if convErr != nil || limit < 0 {
		limit = 0
	}

	tp := u.tokenGenerator.ExtractToken(c)

	usersResponse, err := u.userService.GetAll(
		c.Request().Context(),
		tp.ID,
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

// getOwnProfile godoc
// @Summary      Get Own Profile
// @Description  This endpoint is used to get their own user profile
// @Tags         users
// @Produce      json
// @Security     ApiKey
// @Security     ApiKeyAuth
// @Success      200  {object}  profileResponse
// @Failure      401  {object}  echo.HTTPError
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /users/me [get]
func (u *usersController) getMe(c echo.Context) error {
	tp := u.tokenGenerator.ExtractToken(c)

	userResponse, err := u.userService.GetOwn(c.Request().Context(), tp.ID, tp.Username)
	if err != nil {
		return newErrorResponse(err)
	}

	response := model.NewResponse("success", "Get user successful.", userResponse)

	return c.JSON(http.StatusOK, response)
}

// getUserByUsername godoc
// @Summary      Get User by Username
// @Description  This endpoint is used to get the another user by username
// @Tags         users
// @Produce      json
// @Param        username  path  string  true  "username"
// @Security     ApiKey
// @Security     ApiKeyAuth
// @Success      200  {object}  profileResponse
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /users/{username} [get]
func (u *usersController) getUserByUsername(c echo.Context) error {
	username := c.Param("username")

	tp := u.tokenGenerator.ExtractToken(c)

	userResponse, err := u.userService.GetByUsername(c.Request().Context(), tp.ID, username)
	if err != nil {
		return newErrorResponse(err)
	}

	response := model.NewResponse("success", "Get user successful.", userResponse)

	return c.JSON(http.StatusOK, response)
}

// getUserThreads godoc
// @Summary      Get User Threads
// @Description  This endpoint is used to get the user threads
// @Tags         users
// @Produce      json
// @Param        username  path   string  true   "username"
// @Param        page      query  int     false  "page, default 1"
// @Param        limit     query  int     false  "limit, default 10"
// @Security     ApiKey
// @Security     ApiKeyAuth
// @Success      200  {object}  threadsResponse
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /users/{username}/threads [get]
func (u *usersController) getUserThreads(c echo.Context) error {
	username := c.Param("username")

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

	tp := u.tokenGenerator.ExtractToken(c)

	threadsResponse, err := u.userService.GetAllThreadByUsername(
		c.Request().Context(),
		tp.ID,
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

// putUserFollow godoc
// @Summary      Follow/Unfollow a User
// @Description  This endpoint is used to follow/unfollow a user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        username  path  string  true  "username"
// @Security     ApiKey
// @Security     ApiKeyAuth
// @Success      204
// @Failure      401  {object}  echo.HTTPError
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /users/{username}/follow [put]
func (u *usersController) putUserFollow(c echo.Context) error {
	username := c.Param("username")

	tp := u.tokenGenerator.ExtractToken(c)

	if err := u.userService.ChangeFollowingState(
		c.Request().Context(),
		tp.ID,
		username,
	); err != nil {
		return newErrorResponse(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// putUserBanned godoc
// @Summary      Banned/Unbanned a User
// @Description  This endpoint is used to banned/unbanned a user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        username  path  string  true  "username"
// @Security     ApiKey
// @Security     ApiKeyAuth
// @Success      204
// @Failure      400  {object}  echo.HTTPError
// @Failure      401  {object}  echo.HTTPError
// @Failure      403  {object}  echo.HTTPError
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /users/{username}/banned [put]
func (u *usersController) putUserBanned(c echo.Context) error {
	username := c.Param("username")

	tp := u.tokenGenerator.ExtractToken(c)

	if err := u.userService.ChangeBannedState(
		c.Request().Context(),
		tp.Role,
		username,
	); err != nil {
		return newErrorResponse(err)
	}

	return c.NoContent(http.StatusNoContent)
}

// profileResponse struct is used for swaggo to generate the API documentation, as it doesn't support generic yet.
type profileResponse struct {
	Status  string        `json:"status" extensions:"x-order=0"`
	Message string        `json:"message" extensions:"x-order=1"`
	Data    response.User `json:"data" extensions:"x-order=2"`
}

// profilesResponse struct is used for swaggo to generate the API documentation, as it doesn't support generic yet.
type profilesResponse struct {
	Status  string              `json:"status" extensions:"x-order=0"`
	Message string              `json:"message" extensions:"x-order=1"`
	Data    profilesInfoWrapper `json:"data" extensions:"x-order=2"`
}

// threadsResponse struct is used for swaggo to generate the API documentation, as it doesn't support generic yet.
type threadsResponse struct {
	Status  string             `json:"status" extensions:"x-order=0"`
	Message string             `json:"message" extensions:"x-order=1"`
	Data    threadsInfoWrapper `json:"data" extensions:"x-order=2"`
}

type profilesInfoWrapper struct {
	Threads  []response.User `json:"list" extensions:"x-order=0"`
	PageInfo pageInfoData    `json:"pageInfo" extensions:"x-order=1"`
}

type threadsInfoWrapper struct {
	Threads  []response.ManyThread `json:"list" extensions:"x-order=0"`
	PageInfo pageInfoData          `json:"pageInfo" extensions:"x-order=1"`
}

type moderatorData struct {
	ModeratorID string `json:"moderatorID" extensions:"x-order=0"`
	UserID      string `json:"userID" extensions:"x-order=1"`
	Username    string `json:"username" extensions:"x-order=2"`
	Email       string `json:"email" extensions:"x-order=3"`
	Name        string `json:"name" extensions:"x-order=4"`
	Role        string `json:"role" extensions:"x-order=5"`
	IsActive    bool   `json:"isActive" extensions:"x-order=6"`
	// RegisteredOn layout format: time.RFC822 (02 Jan 06 15:04 MST)
	RegisteredOn string `json:"registeredOn" extensions:"x-order=7"`
}

type pageInfoData struct {
	Limit     uint `json:"limit" extensions:"x-order=0"`
	Page      uint `json:"page" extensions:"x-order=1"`
	PageTotal uint `json:"pageTotal" extensions:"x-order=2"`
	Total     uint `json:"total" extensions:"x-order=3"`
}
