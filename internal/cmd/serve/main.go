package serve

import (
	"github.com/snapp-incubator/nats-readiness/internal/infra/config"
	"github.com/snapp-incubator/nats-readiness/internal/infra/http"
	"github.com/snapp-incubator/nats-readiness/internal/infra/logger"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func main(
	logger *zap.Logger,
	_ http.HTTP,
) {
	logger.Info("welcome to our server")
}

// Register server command.
func Register(
	root *cobra.Command,
) {
	root.AddCommand(
		//nolint: exhaustruct
		&cobra.Command{
			Use:   "serve",
			Short: "Run server to serve the requests",
			Run: func(_ *cobra.Command, _ []string) {
				fx.New(
					fx.Provide(config.Provide),
					fx.Provide(logger.Provide),
					fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
						return &fxevent.ZapLogger{Logger: logger}
					}),
					fx.Provide(http.Provide),
					fx.Invoke(main),
				).Run()
			},
		},
	)
}
