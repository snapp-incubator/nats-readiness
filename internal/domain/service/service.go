package service

import (
	"time"

	"github.com/snapp-incubator/nats-readiness/internal/infra/nats"
)

// NATS service using the connection to the NATS node to figure out the topology
// and the health of the nodes.
type NATS struct {
	connz nats.NATS

	quit chan struct{}
}

func Provide(connz nats.NATS) NATS {
	return NATS{
		connz: connz,
	}
}

// Start run a go routine to calls nats nodes every one minute and gather their information.
func (n NATS) Start() {
	go func() {
		ticker := time.NewTicker(time.Minute)

		for {
			select {
			case <-ticker.C:
				n.connz.Raftz()
			case <-n.quit:
				ticker.Stop()
				return
			}
		}
	}()
}
