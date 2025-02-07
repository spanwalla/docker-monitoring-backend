package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
	"github.com/spanwalla/docker-monitoring-backend/internal/service"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
	"os"
)

// ConfigureRouter -.
func ConfigureRouter(handler *echo.Echo, services *service.Services) {
	handler.Use(middleware.CORS())
	handler.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `{"time":"${time_rfc3339_nano}", "method":"${method}","uri":"${uri}", "status":${status},"error":"${error}"}` + "\n",
		Output: setLogsFile(),
	}))
	handler.Use(middleware.Recover())

	handler.GET("/health", func(c echo.Context) error { return c.NoContent(http.StatusOK) })
	handler.GET("/swagger/*", echoSwagger.WrapHandler)

	authMiddleware := &AuthMiddleware{services.Pinger}
	v1 := handler.Group("/api/v1")
	{
		newPingerRoutes(v1.Group("/pingers"), services.Pinger)
		newReportRoutes(v1.Group("/reports"), services.Report, authMiddleware.PingerIdentity)
	}
}

// setLogsFile -.
func setLogsFile() *os.File {
	file, err := os.OpenFile("/logs/requests.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("v1 - setLogsFile - os.OpenFile: %v", err)
	}
	return file
}
