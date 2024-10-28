package http

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Config struct {
	Listen string `json:"listen,omitempty" koanf:"listen"`
}

type HTTP struct {
	r *http.ServeMux
}

func Provide(
	lc fx.Lifecycle,
	cfg Config,
	logger *zap.Logger,
) HTTP {
	router := http.NewServeMux()

	router.Handle("GET /metrics", promhttp.Handler())

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Info("starting HTTP server")

			go func() {
				if err := http.ListenAndServe("0.0.0.0:1373", router); err != nil {
					logger.Error("http server failed", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(context.Context) error {
			logger.Info("stopping http server.")
			return nil
		},
	})

	return HTTP{
		r: router,
	}
}
