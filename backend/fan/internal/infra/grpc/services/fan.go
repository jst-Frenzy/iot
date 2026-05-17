package services

import (
	"context"
	pb "github.com/jst-Frenzy/iot/backend/fan/api/grpc/gen"
	"google.golang.org/grpc"
	"log/slog"
)

type Fan struct {
	pb.UnimplementedFanServiceServer
}

func New() *Fan {
	return &Fan{}
}

func (f *Fan) Register(server *grpc.Server) {
	pb.RegisterFanServiceServer(server, f)
}

func (f *Fan) On(ctx context.Context, request *pb.OnRequest) (*pb.OnResponse, error) {
	slog.Info("On")
	return &pb.OnResponse{}, nil
}

func (f *Fan) Off(ctx context.Context, request *pb.OffRequest) (*pb.OffResponse, error) {
	slog.Info("Off")
	return &pb.OffResponse{}, nil
}
