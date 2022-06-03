package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type usersController struct{}

func NewUsersController() *usersController {
	return &usersController{}
}

func (u *usersController) Route(g *echo.Group) {
	group := g.Group("/users")
	group.GET("/", u.getUsers)
	group.GET("/me", u.getMe)
	group.GET("/:username", u.getUserByUsername)
	group.GET("/:username/threads", u.getUserThreads)
	group.PUT("/:username/follow", u.putUserFollow)
	group.PUT("/:username/banned", u.putUserBanned)
}

// getUsers     godoc
// @Summary      Get Users
// @Description  This endpoint is used to get all users
// @Tags         users
// @Produce      json
// @Param        page      query     int     false  "page, default 1"
// @Param        limit     query     int     false  "limit, default 20"
// @Param        order_by  query     string  false  "options: registered_date, ranking, default registered_date"
// @Param        status    query     string  false  "options: active, banned, default active"
// @Success      200       {object}  profilesResponse
// @Failure      500       {object}  echo.HTTPError
// @Router       /users [get]
func (u *usersController) getUsers(c echo.Context) error {
	_ = c.QueryParam("page")
	_ = c.QueryParam("limit")
	_ = c.QueryParam("order_by")
	_ = c.QueryParam("status")
	return nil
}

// getOwnProfile godoc
// @Summary      Get Own Profile
// @Description  This endpoint is used to get their own user profile
// @Tags         users
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200       {object}  profileResponse
// @Failure      401  {object}  echo.HTTPError
// @Failure      404       {object}  echo.HTTPError
// @Failure      500       {object}  echo.HTTPError
// @Router       /users/me [get]
func (u *usersController) getMe(c echo.Context) error {
	return nil
}

// getUserByUsername godoc
// @Summary      Get User by Username
// @Description  This endpoint is used to get the another user by username
// @Tags         users
// @Produce      json
// @Param        username  path      string  true  "username"
// @Success      200  {object}  profileResponse
// @Failure      404       {object}  echo.HTTPError
// @Failure      500       {object}  echo.HTTPError
// @Router       /users/{username} [get]
func (u *usersController) getUserByUsername(c echo.Context) error {
	_ = c.Param("username")
	return nil
}

// getUserThreads godoc
// @Summary      Get User Threads
// @Description  This endpoint is used to get the user threads
// @Tags         users
// @Produce      json
// @Param        username  path      string  true   "username"
// @Param        page      query     int     false  "page, default 1"
// @Param        limit     query     int     false  "limit, default 10"
// @Success      200       {object}  threadsResponse
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /users/{username}/threads [get]
func (u *usersController) getUserThreads(c echo.Context) error {
	_ = c.Param("username")
	_ = c.QueryParam("page")
	_ = c.QueryParam("limit")
	return nil
}

// putUserFollow godoc
// @Summary      Follow/Unfollow a User
// @Description  This endpoint is used to follow/unfollow a user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        username  path  string  true  "username"
// @Security     ApiKeyAuth
// @Success      204
// @Failure      401  {object}  echo.HTTPError
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /users/{username}/follow [put]
func (u *usersController) putUserFollow(c echo.Context) error {
	_ = c.Param("username")
	return c.NoContent(http.StatusNoContent)
}

// putUserBanned godoc
// @Summary      Banned/Unbanned a User
// @Description  This endpoint is used to banned/unbanned a user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        username  path  string  true  "username"
// @Security     ApiKeyAuth
// @Success      204
// @Failure      401  {object}  echo.HTTPError
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /users/{username}/banned [put]
func (u *usersController) putUserBanned(c echo.Context) error {
	_ = c.Param("username")
	return c.NoContent(http.StatusNoContent)
}

// profileResponse struct is used for swaggo to generate the API documentation, as it doesn't support generic yet.
type profileResponse struct {
	Status  string      `json:"status" extensions:"x-order=0"`
	Message string      `json:"message" extensions:"x-order=1"`
	Data    profileData `json:"data" extensions:"x-order=2"`
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

type profileData struct {
	UserID   string `json:"userID" extensions:"x-order=0"`
	Username string `json:"username" extensions:"x-order=1"`
	Email    string `json:"email" extensions:"x-order=2"`
	Name     string `json:"name" extensions:"x-order=3"`
	Role     string `json:"role" extensions:"x-order=4"`
	IsActive bool   `json:"isActive" extensions:"x-order=5"`
	// RegisteredOn layout format: time.RFC822 (02 Jan 06 15:04 MST)
	RegisteredOn  string `json:"registeredOn" extensions:"x-order=6"`
	TotalThread   uint   `json:"totalThread" extensions:"x-order=7"`
	TotalFollower uint   `json:"totalFollower" extensions:"x-order=8"`
	IsFollowed    bool   `json:"isFollowed" extensions:"x-order=9"`
}

type profilesInfoWrapper struct {
	Threads  []profileData `json:"users" extensions:"x-order=0"`
	PageInfo pageInfoData  `json:"pageInfo" extensions:"x-order=1"`
}

type threadsInfoWrapper struct {
	Threads  []threadData `json:"threads" extensions:"x-order=0"`
	PageInfo pageInfoData `json:"pageInfo" extensions:"x-order=1"`
}

type threadData struct {
	ID              string `json:"ID" extensions:"x-order=0"`
	Title           string `json:"title" extensions:"x-order=1"`
	Description     string `json:"description" extensions:"x-order=2"`
	TotalViewer     uint   `json:"totalViewer" extensions:"x-order=3"`
	TotalLike       uint   `json:"totalLike" extensions:"x-order=4"`
	TotalFollower   uint   `json:"totalFollower" extensions:"x-order=5"`
	TotalComment    uint   `json:"totalComment" extensions:"x-order=6"`
	CreatorID       string `json:"creatorID" extensions:"x-order=7"`
	CreatorUsername string `json:"creatorUsername" extensions:"x-order=8"`
	CreatorName     string `json:"creatorName" extensions:"x-order=9"`
	CategoryID      string `json:"categoryID" extensions:"x-order=10"`
	CategoryName    string `json:"categoryName" extensions:"x-order=11"`
	// PublishedOn layout format: time.RFC822 (02 Jan 06 15:04 MST)
	PublishedOn string        `json:"publishedOn" extensions:"x-order=12"`
	IsLiked     bool          `json:"isLiked" extensions:"x-order=13"`
	IsFollowed  bool          `json:"isFollowed" extensions:"x-order=14"`
	Moderators  []profileData `json:"moderators" extensions:"x-order=15"`
}

type pageInfoData struct {
	Limit     uint `json:"limit" extensions:"x-order=0"`
	Page      uint `json:"page" extensions:"x-order=1"`
	PageTotal uint `json:"pageTotal" extensions:"x-order=2"`
	Total     uint `json:"total" extensions:"x-order=3"`
}
