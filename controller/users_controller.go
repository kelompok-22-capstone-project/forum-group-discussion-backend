package controller

import "github.com/labstack/echo/v4"

type usersController struct{}

func NewUsersController() *usersController {
	return &usersController{}
}

func (u *usersController) Route(g *echo.Group) {
	group := g.Group("/users")
	group.GET("/me", u.getMe)
	group.GET("/:username", u.getUserByUsername)
}

// getOwnProfile godoc
// @Summary      Get Own Profile
// @Description  This endpoint is used to get their own user profile
// @Tags         users
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {object}  profileResponse
// @Failure      401  {object}  echo.HTTPError
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /users/me [get]
func (u *usersController) getMe(c echo.Context) error {
	return nil
}

// getUserByUsername godoc
// @Summary      Get User by Username
// @Description  This endpoint is used to get the another user by username
// @Tags         users
// @Produce      json
// @Param        username          path  string  true  "username"
// @Success      200  {object}  profileResponse
// @Failure      401  {object}  echo.HTTPError
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /users/{username} [get]
func (u *usersController) getUserByUsername(c echo.Context) error {
	_ = c.Param("username")
	return nil
}

// profileResponse struct is used for swaggo to generate the API documentation, as it doesn't support generic yet.
type profileResponse struct {
	Status  string      `json:"status" extensions:"x-order=0"`
	Message string      `json:"message" extensions:"x-order=1"`
	Data    profileData `json:"data" extensions:"x-order=2"`
}

type profileData struct {
	UserID   string `json:"userID" extensions:"x-order=0"`
	Username string `json:"username" extensions:"x-order=1"`
	Email    string `json:"email" extensions:"x-order=2"`
	Name     string `json:"name" extensions:"x-order=3"`
	Role     string `json:"role" extensions:"x-order=4"`
	IsActive bool   `json:"isActive" extensions:"x-order=5"`
	// RegisteredOn layout format: time.RFC822 (02 Jan 06 15:04 MST)
	RegisteredOn  string       `json:"registeredOn" extensions:"x-order=6"`
	TotalThread   uint         `json:"totalThread" extensions:"x-order=7"`
	TotalFollower uint         `json:"totalFollower" extensions:"x-order=8"`
	Threads       []threadData `json:"threads" extensions:"x-order=9"`
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
	CategoryID      string `json:"categoryID" extensions:"x-order=9"`
	CategoryName    string `json:"categoryName" extensions:"x-order=10"`
	// PublishedOn layout format: time.RFC822 (02 Jan 06 15:04 MST)
	PublishedOn string `json:"publishedOn" extensions:"x-order=11"`
}
