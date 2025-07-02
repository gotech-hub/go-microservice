package http

import (
	"net/http"

	"go-source/pkg/binding"
	middlewares "go-source/pkg/middlewares"
	"go-source/pkg/resp"

	_ "go-source/docs" // swagger docs

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

const healthPath = "/v1/service-name/health"

func (app *Server) InitRouters(e *echo.Echo) error {
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middlewares.Logging)
	e.Use(middlewares.AddExtraDataForRequestContext)

	// Swagger documentation endpoint
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET(healthPath, func(c echo.Context) error {
		if healthCheck {
			return c.JSON(http.StatusOK, resp.BuildSuccessResp(resp.LangEN, nil))
		}

		return c.JSON(http.StatusInternalServerError, resp.BuildErrorResp(500, "", resp.LangEN))
	})

	e.POST("/v1/service-name/test", binding.Wrapper(app.Handlers.Handler.GetByProfileId))

	return nil
}
