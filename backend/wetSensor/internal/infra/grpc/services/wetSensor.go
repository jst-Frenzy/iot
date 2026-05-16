package services

import (
	"context"
	pb "github.com/jst-Frenzy/iot/backend/wetSensor/api/grpc/gen"
	"google.golang.org/grpc"
)

type device interface {
	GenerateData() int32
}

type WetSensorDeps struct {
	Device device
}

type WetSensor struct {
	pb.UnimplementedWetSensorServiceServer
	device device
}

func New(d *WetSensorDeps) *WetSensor {
	return &WetSensor{
		device: d.Device,
	}
}

func (t *WetSensor) Register(server *grpc.Server) {
	pb.RegisterWetSensorServiceServer(server, t)
}

func (t *WetSensor) GetData(ctx context.Context, request *pb.GetDataRequest) (*pb.GetDataResponse, error) {
	return &pb.GetDataResponse{
		Data: t.device.GenerateData(),
	}, nil
}
