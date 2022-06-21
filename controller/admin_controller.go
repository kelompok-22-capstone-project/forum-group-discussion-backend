package controller

import (
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/response"
	"github.com/labstack/echo/v4"
)

type adminController struct{}

func NewAdminController() *adminController {
	return &adminController{}
}

func (i *adminController) Route(g *echo.Group) {
	group := g.Group("/admin/dashboard")
	group.GET("/", i.getInfo)
}

// getInfo     godoc
// @Summary      Get Info
// @Description  This endpoint is used to get all information for admin dashboard purpose
// @Tags         admin
// @Produce      json
// @Security     ApiKey
// @Security     ApiKeyAuth
// @Success      200  {object}  infoResponse
// @Failure      500  {object}  echo.HTTPError
// @Router       /admin/dashboard [get]
func (i *adminController) getInfo(c echo.Context) error {
	return nil
}

// profileResponse struct is used for swaggo to generate the API documentation, as it doesn't support generic yet.
type infoResponse struct {
	Status  string                 `json:"status" extensions:"x-order=0"`
	Message string                 `json:"message" extensions:"x-order=1"`
	Data    response.DashboardInfo `json:"data" extensions:"x-order=2"`
}
