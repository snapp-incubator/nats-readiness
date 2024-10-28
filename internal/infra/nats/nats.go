package nats

import (
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type Config struct {
	Endpoints []string
}

type NATS struct {
	clients []*resty.Client
	logger  *zap.Logger
}

func Provide(logger *zap.Logger, cfg Config) NATS {
	clients := make([]*resty.Client, 0)

	for _, endpoint := range cfg.Endpoints {
		clients = append(clients, resty.New().SetBaseURL(endpoint))
	}

	return NATS{
		clients: clients,
		logger:  logger,
	}
}

func (n NATS) Raftz() {
	for _, client := range n.clients {
		resp, err := client.R().SetQueryParam("js-enabled-only", "1").Get("/healthz")
		if err != nil {
			n.logger.Error("failed to call nats healthz endpoint", zap.Error(err), zap.String("url", client.BaseURL))
		}

		n.logger.Info("nats healthz response", zap.ByteString("response", resp.Body()))
	}
}
