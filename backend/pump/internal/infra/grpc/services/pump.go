package services

import (
	"context"
	pb "github.com/jst-Frenzy/iot/backend/pump/api/grpc/gen"
	"google.golang.org/grpc"
	"log/slog"
)

type Pump struct {
	pb.UnimplementedPumpServiceServer
}

func New() *Pump {
	return &Pump{}
}

func (f *Pump) Register(server *grpc.Server) {
	pb.RegisterPumpServiceServer(server, f)
}

func (f *Pump) On(ctx context.Context, request *pb.OnRequest) (*pb.OnResponse, error) {
	slog.Info("On")
	return &pb.OnResponse{}, nil
}

func (f *Pump) Off(ctx context.Context, request *pb.OffRequest) (*pb.OffResponse, error) {
	slog.Info("Off")
	return &pb.OffResponse{}, nil
}
