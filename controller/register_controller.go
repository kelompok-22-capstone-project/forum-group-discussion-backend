package controller

import (
	"net/http"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/payload"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service/user"
	"github.com/labstack/echo/v4"
)

type registerController struct {
	userService user.UserService
}

func NewRegisterController(userService user.UserService) *registerController {
	return &registerController{userService: userService}
}

func (r *registerController) Route(g *echo.Group) {
	group := g.Group("/register")
	group.POST("", r.postRegister)
}

// PostRegister  godoc
// @Summary      User Register
// @Description  This endpoint is used for user register.
// @Tags         register
// @Accept       json
// @Security     ApiKey
// @Produce      json
// @Security     ApiKey
// @Param        default  body      payload.Register  true  "register payload"
// @Success      201      {object}  registerResponse
// @Failure      400      {object}  echo.HTTPError
// @Failure      500      {object}  echo.HTTPError
// @Router       /register [post]
func (r *registerController) postRegister(c echo.Context) error {
	registerPayload := new(payload.Register)
	if err := c.Bind(registerPayload); err != nil {
		return newErrorResponse(service.ErrInvalidPayload)
	}

	id, err := r.userService.Register(c.Request().Context(), *registerPayload)
	if err != nil {
		return newErrorResponse(err)
	}

	idResponse := map[string]any{"userID": id}
	response := model.NewResponse("success", "Register successful.", idResponse)
	return c.JSON(http.StatusCreated, response)
}

// registerResponse struct is used for swaggo to generate the API documentation, as it doesn't support generic yet.
type registerResponse struct {
	Status  string     `json:"status" extensions:"x-order=0"`
	Message string     `json:"message" extensions:"x-order=1"`
	Data    userIDData `json:"data" extensions:"x-order=2"`
}

type userIDData struct {
	UserID string `json:"userID"`
}
