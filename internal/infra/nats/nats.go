package nats

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type HTTPError struct {
	endpoint  string
	error     error
	status    int
	operation string
}

func (h HTTPError) Error() string {
	if h.error != nil {
		return fmt.Sprintf("failed to call %s endpoint on %s, %s", h.operation, h.endpoint, h.error)
	}

	return fmt.Sprintf("failed to call %s endpoint on %s, status %d", h.operation, h.endpoint, h.status)
}

func (h HTTPError) Unwrap() error {
	if h.error != nil {
		return h.error
	}

	return nil
}

type Raftz struct {
	SYS struct {
		Meta struct {
			ID             string `json:"id"`
			State          string `json:"state"`
			Size           int    `json:"size"`
			QuorumNeeded   int    `json:"quorum_needed"`
			Committed      int    `json:"committed"`
			Applied        int    `json:"applied"`
			Leader         string `json:"leader"`
			EverHadLeader  bool   `json:"ever_had_leader"`
			Term           int    `json:"term"`
			VotedFor       string `json:"voted_for"`
			Pterm          int    `json:"pterm"`
			Pindex         int    `json:"pindex"`
			IpqProposalLen int    `json:"ipq_proposal_len"`
			IpqEntryLen    int    `json:"ipq_entry_len"`
			IpqRespLen     int    `json:"ipq_resp_len"`
			IpqApplyLen    int    `json:"ipq_apply_len"`
			Wal            struct {
				Messages      int       `json:"messages"`
				Bytes         int       `json:"bytes"`
				FirstSeq      int       `json:"first_seq"`
				FirstTs       time.Time `json:"first_ts"`
				LastSeq       int       `json:"last_seq"`
				LastTs        time.Time `json:"last_ts"`
				ConsumerCount int       `json:"consumer_count"`
			} `json:"wal"`
			Peers map[string]struct {
				Name     string `json:"name"`
				Known    bool   `json:"known"`
				LastSeen string `json:"last_seen"`
			} `json:"peers"`
		} `json:"_meta_"`
	} `json:"$SYS"`
}

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

func (n NATS) Healthz() error {
	errs := make([]error, 0)

	for _, client := range n.clients {
		resp, err := client.R().SetQueryParam("js-enabled-only", "1").Get("/healthz")
		if err != nil {
			n.logger.Error("failed to call nats healthz endpoint", zap.Error(err), zap.String("url", client.BaseURL))

			errs = append(errs, HTTPError{
				endpoint:  client.BaseURL,
				error:     err,
				status:    -1,
				operation: "healthz",
			})

			continue
		}

		n.logger.Debug("nats healthz response", zap.ByteString("response", resp.Body()))

		if !resp.IsSuccess() {
			errs = append(errs, HTTPError{
				endpoint:  client.BaseURL,
				error:     nil,
				status:    resp.StatusCode(),
				operation: "healthz",
			})

			continue
		}

	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (n NATS) Raftz() {
	for _, client := range n.clients {
		resp, err := client.R().Get("/raftz")
		if err != nil {
			n.logger.Error("failed to call nats raftz endpoint", zap.Error(err), zap.String("url", client.BaseURL))
		}

		n.logger.Debug("nats healthz response", zap.ByteString("response", resp.Body()))
	}
}
