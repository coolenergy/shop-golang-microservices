package server

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/grpc"
	echoserver "github.com/meysamhadeli/shop-golang-microservices/internal/pkg/http/echo/server"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/identity_service/config"
	"github.com/pkg/errors"
	"go.uber.org/fx"
	"net/http"
)

func RunServers(lc fx.Lifecycle, log logger.ILogger, e *echo.Echo, grpcServer *grpc.GrpcServer, ctx context.Context, cfg *config.Config) error {

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				if err := echoserver.RunHttpServer(ctx, e, log, cfg.Echo); !errors.Is(err, http.ErrServerClosed) {
					log.Fatalf("error running http server: %v", err)
				}
			}()

			go func() {
				if err := grpcServer.RunGrpcServer(ctx); !errors.Is(err, http.ErrServerClosed) {
					log.Fatalf("error running grpc server: %v", err)
				}
			}()

			e.GET("/", func(c echo.Context) error {
				return c.String(http.StatusOK, config.GetMicroserviceName(cfg.ServiceName))
			})

			return nil
		},
		OnStop: func(_ context.Context) error {
			log.Infof("all servers shutdown gracefully...")
			return nil
		},
	})

	return nil
}
