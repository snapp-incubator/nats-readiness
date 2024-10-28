package config

import (
	"github.com/snapp-incubator/nats-readiness/internal/infra/http"
	"github.com/snapp-incubator/nats-readiness/internal/infra/logger"
	"github.com/snapp-incubator/nats-readiness/internal/infra/nats"
)

func Default() Config {
	return Config{
		HTTP: http.Config{
			Listen: "0.0.0.0:1373",
		},
		Logger: logger.Config{
			Level: "debug",
		},
		NATS: nats.Config{
			Endpoints: []string{
				"http://127.0.0.1:8222",
			},
		},
	}
}
