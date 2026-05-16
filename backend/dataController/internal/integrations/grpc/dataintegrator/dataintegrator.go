package dataintegrator

import (
	"context"
	"fmt"
	dic "github.com/jst-Frenzy/iot/backend/dataIntegrator/api/grpc/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
)

type Client struct {
	service dic.DataIntegratorServiceClient
	logger  *slog.Logger
}

type Config struct {
	Address        string
	MaxMessageSize int
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
		service: dic.NewDataIntegratorServiceClient(conn),
		logger:  slog.With(slog.With("service", "integrations.grpc.temperatureService")),
	}, nil
}

func (c *Client) SendData(ctx context.Context, wet, temperature int32) error {
	_, err := c.service.AcceptData(ctx, &dic.AcceptDataRequest{
		Wet:         wet,
		Temperature: temperature,
	})
	return err
}
