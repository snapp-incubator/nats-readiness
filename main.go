package main

import (
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()

	client := resty.New()

	resp, err := client.R().SetQueryParam("js-enabled-only", "1").Get("http://127.0.0.1:8222/healthz")
	if err != nil {
		logger.Error("failed to call nats healthz endpoint", zap.Error(err))
	}

	logger.Info("nats healthz response", zap.ByteString("response", resp.Body()))
}
