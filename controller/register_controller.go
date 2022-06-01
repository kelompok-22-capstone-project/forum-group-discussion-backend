package controller

import (
	"net/http"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/payload"
	"github.com/labstack/echo/v4"
)

type registerController struct{}

func NewRegisterController() *registerController {
	return &registerController{}
}

func (r *registerController) Route(g *echo.Group) {
	group := g.Group("/register")
	group.POST("", r.postRegister)
}

func (p *registerController) postRegister(c echo.Context) error {
	registerPayload := new(payload.Register)
	if err := c.Bind(registerPayload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid payload. Please check the payload schema in the API Documentation.")
	}

	id := "u-xy4Zia"

	idResponse := map[string]any{"id": id}
	response := model.NewResponse("success", "Register successful.", idResponse)
	return c.JSON(http.StatusOK, response)
}
