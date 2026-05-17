package services

import (
	"context"
	pb "github.com/jst-Frenzy/iot/backend/fan/api/grpc/gen"
	"google.golang.org/grpc"
)

type device interface {
	GenerateData() int32
}

type FanDeps struct {
	Device device
}

type Fan struct {
	pb.UnimplementedFanServiceServer
	device device
}

func New(d *FanDeps) *Fan {
	return &Fan{
		device: d.Device,
	}
}

func (f *Fan) Register(server *grpc.Server) {
	pb.RegisterFanServiceServer(server, f)
}

func (f *Fan) On(ctx context.Context, request *pb.OnRequest) (*pb.OnResponse, error) {
	return &pb.OnResponse{}, nil
}

func (f *Fan) Off(ctx context.Context, request *pb.OffRequest) (*pb.OffResponse, error) {
	return &pb.OffResponse{}, nil
}
