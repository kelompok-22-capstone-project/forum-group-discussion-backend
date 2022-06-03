package controller

import "github.com/labstack/echo/v4"

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
// @Security     ApiKeyAuth
// @Success      200  {object}  infoResponse
// @Failure      500  {object}  echo.HTTPError
// @Router       /admin/dashboard [get]
func (i *adminController) getInfo(c echo.Context) error {
	return nil
}

// profileResponse struct is used for swaggo to generate the API documentation, as it doesn't support generic yet.
type infoResponse struct {
	Status  string   `json:"status" extensions:"x-order=0"`
	Message string   `json:"message" extensions:"x-order=1"`
	Data    infoData `json:"data" extensions:"x-order=2"`
}

type infoData struct {
	TotalUser      uint `json:"totalUser" extensions:"x-order=0"`
	TotalThread    uint `json:"totalThread" extensions:"x-order=1"`
	TotalModerator uint `json:"totalModerator" extensions:"x-order=2"`
	TotalReport    uint `json:"totalReport" extensions:"x-order=3"`
}
