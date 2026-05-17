package grpc

import (
	"context"
	"errors"
	"github.com/jst-Frenzy/iot/backend/pump/internal/config/configuration"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"
	"time"
)

type Service interface {
	Register(*grpc.Server)
}

type Server struct {
	srv      *grpc.Server
	listener net.Listener
	address  string
	logger   *slog.Logger
}

type Deps struct {
	Conf     *configuration.GRPCServer `validate:"required"`
	Logger   *slog.Logger              `validate:"required"`
	Services []Service
}

func New(dep *Deps) (*Server, error) {
	srv := grpc.NewServer(
		grpc.MaxRecvMsgSize(dep.Conf.MaxReceiveMsgSize),
		grpc.MaxSendMsgSize(dep.Conf.MaxSendMsgSize),
	)

	for _, svc := range dep.Services {
		svc.Register(srv)
	}

	reflection.Register(srv)

	return &Server{
		srv:     srv,
		address: dep.Conf.Address,
		logger:  dep.Logger,
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	conf := &net.ListenConfig{}

	lis, err := conf.Listen(ctx, "tcp", s.address)
	if err != nil {
		return err
	}

	s.listener = lis

	errChan := make(chan error, 1)
	go func() {
		s.logger.InfoContext(ctx, "gRPC listening", "address", s.address)
		if err := s.srv.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			errChan <- err
		}
		close(errChan)
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		s.logger.InfoContext(ctx, "shutting down gRPC server")

		stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		return s.Stop(stopCtx)
	}
}

func (s *Server) Stop(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		s.srv.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		s.logger.InfoContext(ctx, "force stopping gRPC server")
		s.srv.Stop()
		return ctx.Err()
	}
}
