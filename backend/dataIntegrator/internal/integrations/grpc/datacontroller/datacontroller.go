package datacontroller

import (
	"context"
	"fmt"
	"log/slog"

	dcs "github.com/jst-Frenzy/iot/backend/dataController/api/grpc/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	service dcs.DeviceServiceClient
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
		service: dcs.NewDeviceServiceClient(conn),
		logger:  slog.With(slog.With("service", "integrations.grpc.temperatureService")),
	}, nil
}

func (c *Client) OnFan(ctx context.Context) error {
	_, err := c.service.OnFan(ctx, &dcs.OnFanRequest{})
	return err
}

func (c *Client) OffFan(ctx context.Context) error {
	_, err := c.service.OffFan(ctx, &dcs.OffFanRequest{})
	return err
}

func (c *Client) OnPump(ctx context.Context) error {
	_, err := c.service.OnPump(ctx, &dcs.OnPumpRequest{})
	return err
}

func (c *Client) OffPump(ctx context.Context) error {
	_, err := c.service.OffPump(ctx, &dcs.OffPumpRequest{})
	return err
}
