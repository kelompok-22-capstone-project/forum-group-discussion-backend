package controller

import "github.com/labstack/echo/v4"

type reportsController struct{}

func NewReportsController() *reportsController {
	return &reportsController{}
}

func (r *reportsController) Route(g *echo.Group) {
	group := g.Group("/reports")
	group.POST("", r.postCreateReport)
	group.GET("", r.getReports)
	group.PUT("/:id/status", r.putUpdateReportStatus)
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

// putUpdateReportStatus godoc
// @Summary      Update a Report Status
// @Description  This endpoint is used to update a report status
// @Tags         reports
// @Accept       json
// @Produce      json
// @Param        default  body  payload.UpdateReportStatus  true  "request body"
// @Param        id       path  string                      true  "report ID"
// @Security     ApiKey
// @Security     ApiKeyAuth
// @Success      204
// @Failure      400  {object}  echo.HTTPError
// @Failure      401  {object}  echo.HTTPError
// @Failure      404  {object}  echo.HTTPError
// @Failure      500  {object}  echo.HTTPError
// @Router       /reports/{id}/status [put]
func (r *reportsController) putUpdateReportStatus(c echo.Context) error {
	_ = c.Param("id")
	return nil
}

// reportsResponse struct is used for swaggo to generate the API documentation, as it doesn't support generic yet.
type reportsResponse struct {
	Status  string      `json:"status" extensions:"x-order=0"`
	Message string      `json:"message" extensions:"x-order=1"`
	Data    reportsData `json:"data" extensions:"x-order=2"`
}

type reportsData struct {
	Reports  []reportData `json:"list" extensions:"x-order=0"`
	PageInfo pageInfoData `json:"pageInfo" extensions:"x-order=1"`
}

type reportData struct {
	ID                string `json:"ID" extensions:"x-order=0"`
	ModeratorID       string `json:"moderatorID" extensions:"x-order=1"`
	ModeratorUsername string `json:"moderatorUsername" extensions:"x-order=2"`
	ModeratorName     string `json:"moderatorName" extensions:"x-order=3"`
	UserID            string `json:"userID" extensions:"x-order=4"`
	Username          string `json:"username" extensions:"x-order=5"`
	Name              string `json:"name" extensions:"x-order=6"`
	Reason            string `json:"reason" extensions:"x-order=7"`
	Status            string `json:"status" extensions:"x-order=8"`
	// ReportedOn layout format: time.RFC822 (02 Jan 06 15:04 MST)
	ReportedOn string `json:"reportedOn" extensions:"x-order=9"`
}
