package controller

import (
	"net/http"
	"strconv"

	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/middleware"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/payload"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/model/response"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/service/report"
	"github.com/kelompok-22-capstone-project/forum-group-discussion-backend/utils/generator"
	"github.com/labstack/echo/v4"
)

type reportsController struct {
	service        report.ReportService
	tokenGenerator generator.TokenGenerator
}

func NewReportsController(
	service report.ReportService,
	tokenGenerator generator.TokenGenerator,
) *reportsController {
	return &reportsController{service: service, tokenGenerator: tokenGenerator}
}

func (r *reportsController) Route(g *echo.Group) {
	group := g.Group("/reports")
	group.POST("", r.postCreateReport, middleware.JWTMiddleware())
	group.GET("", r.getReports, middleware.JWTMiddleware())
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
	p := new(payload.CreateReport)
	if err := c.Bind(p); err != nil {
		return newErrorResponse(service.ErrInvalidPayload)
	}

	tp := r.tokenGenerator.ExtractToken(c)

	id, err := r.service.Create(c.Request().Context(), tp.ID, *p)
	if err != nil {
		return newErrorResponse(err)
	}

	idResponse := map[string]any{"ID": id}
	response := model.NewResponse("success", "Create report successful.", idResponse)
	return c.JSON(http.StatusCreated, response)
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
	status := c.QueryParam("status")
	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")

	page, convErr := strconv.Atoi(pageStr)
	if convErr != nil {
		page = 0
	}

	limit, convErr := strconv.Atoi(limitStr)
	if convErr != nil {
		limit = 0
	}

	tp := r.tokenGenerator.ExtractToken(c)

	reportsResponse, err := r.service.GetAll(c.Request().Context(), tp.Role, status, uint(page), uint(limit))
	if err != nil {
		return newErrorResponse(err)
	}

	response := model.NewResponse("success", "Get reports successful.", reportsResponse)
	return c.JSON(http.StatusOK, response)
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
