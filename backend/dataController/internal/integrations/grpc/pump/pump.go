package pump

import (
	"fmt"
	fc "github.com/jst-Frenzy/iot/backend/fan/api/grpc/gen"
	pc "github.com/jst-Frenzy/iot/backend/pump/api/grpc/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	service pc.
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
		service: fc.NewFanServiceClient(conn),
		logger:  slog.With(slog.With("service", "integrations.grpc.temperatureService")),
	}, nil
}

func (c *Client) On(ctx context.Context) error {
	_, err := c.service.On(ctx, &fc.OnRequest{})
	return err
}

func (c *Client) Off(ctx context.Context) error {
	_, err := c.service.Off(ctx, &fc.OffRequest{})
	return err
}
