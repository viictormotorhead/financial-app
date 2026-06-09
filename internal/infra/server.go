package infra

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/viictormotorhead/financial-app/internal/infra/api/response"
	"github.com/viictormotorhead/financial-app/internal/infra/config"
)

func NewEchoServer(lc fx.Lifecycle, logger *zap.Logger) *echo.Echo {
	e := echo.New()
	e.HTTPErrorHandler = response.HTTPErrorHandler
	e.HideBanner = true
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPatch,
			http.MethodOptions,
		},
	}))

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			addr := fmt.Sprintf(":%s", config.Config().Server.Port)
			logger.Info("starting HTTP server", zap.String("addr", addr))

			go func() {
				if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
					logger.Error("HTTP server failed", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("shutting down HTTP server")
			return e.Shutdown(ctx)
		},
	})

	return e
}
