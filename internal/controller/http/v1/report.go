package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/spanwalla/docker-monitoring-backend/internal/service"
	"net/http"
)

// reportRoutes -.
type reportRoutes struct {
	reportService service.Report
}

// newPingerRoutes -.
func newReportRoutes(g *echo.Group, reportService service.Report, authMiddleware echo.MiddlewareFunc) {
	r := &reportRoutes{reportService}

	g.GET("", r.getReports)
	g.POST("", r.storeReport, authMiddleware)
}

// @Summary Get reports
// @Description Get latest report by every pinger ever exists in database
// @Tags reports
// @Produce json
// @Success 200 {object} v1.reportRoutes.getReports.response
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /api/v1/reports [get]
func (r *reportRoutes) getReports(c echo.Context) error {
	reports, err := r.reportService.GetActualReports(c.Request().Context())
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	type response struct {
		Reports []service.ReportOutput `json:"reports"`
	}

	return c.JSON(http.StatusOK, response{reports})
}

type storeReportInput struct {
	Report []service.PingResult `json:"report" validate:"required"`
}

// @Summary Store report
// @Description Store pinger's report to database
// @Tags reports
// @Accept json
// @Produce json
// @Param input body storeReportInput true "input"
// @Success 202
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Security JWT
// @Router /api/v1/reports [post]
func (r *reportRoutes) storeReport(c echo.Context) error {
	var input storeReportInput
	if err := c.Bind(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	err := r.reportService.PublishToQueue(c.Request().Context(), service.ReportStoreInput{
		PingerId: c.Get(pingerIdCtx).(int),
		Report:   input.Report,
	})
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	return c.JSON(http.StatusAccepted, map[string]string{
		"message": "accepted",
	})
}
