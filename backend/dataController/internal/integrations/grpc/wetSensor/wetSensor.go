package temperatureSensor

import (
	"context"
	"fmt"
	wsc "github.com/jst-Frenzy/iot/backend/wetSensor/api/grpc/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"time"
)

type Client struct {
	service wsc.WetSensorServiceClient
	logger  *slog.Logger
}

type Config struct {
	Address        string
	MaxMessageSize int
	RetryAttempts  int
	MinRetryTime   time.Duration
	MaxRetryTime   time.Duration
}

func New(d *Config) (*Client, error) {
	conn, err := grpc.NewClient(d.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(d.MaxMessageSize)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	return &Client{
		service: wsc.NewWetSensorServiceClient(conn),
		logger:  slog.With(slog.With("service", "integrations.grpc.temperatureService")),
	}, nil
}

func (c *Client) GetWet(ctx context.Context) (int32, error) {
	resp, err := c.service.GetData(ctx, &wsc.GetDataRequest{})
	if err != nil {
		return 0, fmt.Errorf("cant get wet: %w", err)
	}
	return resp.Data, nil
}
