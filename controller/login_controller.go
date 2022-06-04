package controller

import (
	"net/http"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/payload"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/response"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service/user"
	"github.com/labstack/echo/v4"
)

type loginController struct {
	userService user.UserService
}

func NewLoginController(userService user.UserService) *loginController {
	return &loginController{userService: userService}
}

func (l *loginController) Route(g *echo.Group) {
	group := g.Group("/login")
	group.POST("", l.postLogin)
}

// PostLogin     godoc
// @Summary      User Login
// @Description  This endpoint is used for user login
// @Tags         login
// @Accept       json
// @Produce      json
// @Param        default  body      payload.Login  true  "user credentials"
// @Success      200      {object}  loginResponse
// @Failure      400      {object}  echo.HTTPError
// @Failure      401      {object}  echo.HTTPError
// @Failure      404      {object}  echo.HTTPError
// @Failure      500      {object}  echo.HTTPError
// @Router       /login [post]
func (l *loginController) postLogin(c echo.Context) error {
	credential := new(payload.Login)
	if err := c.Bind(credential); err != nil {
		return newErrorResponse(service.ErrInvalidPayload)
	}

	tokenResponse, err := l.userService.Login(c.Request().Context(), *credential)
	if err != nil {
		return newErrorResponse(err)
	}

	response := model.NewResponse("success", "Login successful.", tokenResponse)
	return c.JSON(http.StatusOK, response)
}

// loginResponse struct is used for swaggo to generate the API documentation, as it doesn't support generic yet.
type loginResponse struct {
	Status  string         `json:"status" extensions:"x-order=0"`
	Message string         `json:"message" extensions:"x-order=1"`
	Data    response.Login `json:"data" extensions:"x-order=2"`
}
