package pump

import (
	"context"
	"fmt"
	pc "github.com/jst-Frenzy/iot/backend/pump/api/grpc/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
)

type Client struct {
	service pc.PumpServiceClient
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
		service: pc.NewPumpServiceClient(conn),
		logger:  slog.With(slog.With("service", "integrations.grpc.temperatureService")),
	}, nil
}

func (c *Client) On(ctx context.Context) error {
	_, err := c.service.On(ctx, &pc.OnRequest{})
	return err
}

func (c *Client) Off(ctx context.Context) error {
	_, err := c.service.Off(ctx, &pc.OffRequest{})
	return err
}
