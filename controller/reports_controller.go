package controller

import (
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/response"
	"github.com/labstack/echo/v4"
)

type reportsController struct{}

func NewReportsController() *reportsController {
	return &reportsController{}
}

func (r *reportsController) Route(g *echo.Group) {
	group := g.Group("/reports")
	group.POST("", r.postCreateReport)
	group.GET("", r.getReports)
}

// postCreateReport godoc
// @Summary      Create a Report
// @Description  This endpoint is used to create a report
// @Tags         reports
// @Accept       json
// @Produce      json
// @Param        default  body  payload.CreateReport  true  "request body"
// @Security     ApiKey
// @Security     ApiKeyAuth
// @Success      201  {object}  createThreadResponse
// @Failure      400  {object}  echo.HTTPError
// @Failure      401  {object}  echo.HTTPError
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /reports [post]
func (r *reportsController) postCreateReport(c echo.Context) error {
	return nil
}

// getReports     godoc
// @Summary      Get Reports
// @Description  This endpoint is used to get all report
// @Tags         reports
// @Produce      json
// @Param        status  query  string  false  "options: review, accepted, default review"
// @Param        page    query  int     false  "page, default 1"
// @Param        limit   query  int     false  "limit, default 20"
// @Security     ApiKey
// @Security     ApiKeyAuth
// @Success      200  {object}  reportsResponse
// @Failure      500  {object}  echo.HTTPError
// @Router       /reports [get]
func (r *reportsController) getReports(c echo.Context) error {
	_ = c.QueryParam("status")
	_ = c.QueryParam("page")
	_ = c.QueryParam("limit")
	return nil
}

// reportsResponse struct is used for swaggo to generate the API documentation, as it doesn't support generic yet.
type reportsResponse struct {
	Status  string      `json:"status" extensions:"x-order=0"`
	Message string      `json:"message" extensions:"x-order=1"`
	Data    reportsData `json:"data" extensions:"x-order=2"`
}

type reportsData struct {
	Reports  []response.Report `json:"list" extensions:"x-order=0"`
	PageInfo pageInfoData      `json:"pageInfo" extensions:"x-order=1"`
}
