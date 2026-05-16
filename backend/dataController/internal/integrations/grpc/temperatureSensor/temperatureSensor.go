package temperatureSensor

import (
	"fmt"
	tsc "github.com/jst-Frenzy/iot/backend/temperatureSensor/api/grpc/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
	"time"
)

type Client struct {
	service tsc.DeviceServiceClient
	logger  *slog.Logger
}

type Config struct {
	Address        string        `validate:"required"`
	MaxMessageSize int           `validate:"required"`
	RetryAttempts  int           `validate:"required"`
	MinRetryTime   time.Duration `validate:"required"`
	MaxRetryTime   time.Duration `validate:"required"`
}

func New(d *Config) (*Client, error) {
	if err := validate.Struct(d); err != nil {
		return nil, err
	}

	conn, err := grpc.NewClient(d.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(d.MaxMessageSize)),
		grpc.WithUnaryInterceptor(
			interceptors.UnaryRetry(d.RetryAttempts, d.MinRetryTime, d.MaxRetryTime),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	return &Client{
		service: pbs.NewAuctionServiceClient(conn),
		logger:  slog.With(sl.Component("integrations.grpc.atracs")),
	}, nil
}
