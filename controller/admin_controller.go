package controller

import (
	"net/http"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/middleware"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/response"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service/admin"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator"
	"github.com/labstack/echo/v4"
)

type adminController struct {
	service        admin.AdminService
	tokenGenerator generator.TokenGenerator
}

func NewAdminController(
	service admin.AdminService,
	tokenGenerator generator.TokenGenerator,
) *adminController {
	return &adminController{service: service, tokenGenerator: tokenGenerator}
}

func (i *adminController) Route(g *echo.Group) {
	group := g.Group("/admin/dashboard")
	group.GET("", i.getInfo, middleware.JWTMiddleware())
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
	tp := i.tokenGenerator.ExtractToken(c)

	infoResponse, err := i.service.GetDashboardInfo(c.Request().Context(), tp.Role)
	if err != nil {
		return newErrorResponse(err)
	}

	response := model.NewResponse("Success", "Get admin dashboard info successful.", infoResponse)

	return c.JSON(http.StatusOK, response)
}

// profileResponse struct is used for swaggo to generate the API documentation, as it doesn't support generic yet.
type infoResponse struct {
	Status  string                 `json:"status" extensions:"x-order=0"`
	Message string                 `json:"message" extensions:"x-order=1"`
	Data    response.DashboardInfo `json:"data" extensions:"x-order=2"`
}
