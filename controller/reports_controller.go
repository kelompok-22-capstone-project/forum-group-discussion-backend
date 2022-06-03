package controller

import "github.com/labstack/echo/v4"

type reportsController struct{}

func NewReportsController() *reportsController {
	return &reportsController{}
}

func (r *reportsController) Route(g *echo.Group) {
	group := g.Group("/reports")
	group.POST("", r.postCreateReport)
}

// postCreateReport godoc
// @Summary      Create a Report
// @Description  This endpoint is used to create a report
// @Tags         reports
// @Accept       json
// @Produce      json
// @Param        default  body  payload.CreateReport  true  "request body"
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
