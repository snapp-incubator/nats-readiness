package service

import (
	"context"
	"time"

	"github.com/snapp-incubator/nats-readiness/internal/infra/nats"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// NATS service using the connection to the NATS node to figure out the topology
// and the health of the nodes.
type NATS struct {
	connz  nats.NATS
	logger *zap.Logger

	quit chan struct{}
}

func ProvideNATS(lc fx.Lifecycle, connz nats.NATS, logger *zap.Logger) NATS {
	s := NATS{
		connz:  connz,
		logger: logger.Named("service.nats"),
	}

	lc.Append(fx.Hook{
		OnStart: s.Start,
		OnStop:  s.Stop,
	})

	return s
}

// Start run a go routine to calls nats nodes every one minute and gather their information.
func (n NATS) Start(_ context.Context) error {
	go func() {
		ticker := time.NewTicker(time.Minute)

		for {
			select {
			case <-ticker.C:
				r, err := n.connz.Raftz()
				if err != nil {
					n.logger.Error("failed to call raftz on some endpoints", zap.Error(err))
				}
				n.Update(r)

			case <-n.quit:
				ticker.Stop()
				return
			}
		}
	}()

	return nil
}

func (n NATS) Stop(_ context.Context) error {
	close(n.quit)

	return nil
}

// Update updates the nats cluster view by using the raftz results.
func (n *NATS) Update(r map[string]nats.Raftz) {
	for e, r := range r {
		n.logger.Info("node cluster view",
			zap.String("ip", e),
			zap.String("id", r.SYS.Meta.ID),
		)
	}
}
