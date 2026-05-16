package services

import (
	"context"
	pb "github.com/jst-Frenzy/iot/backend/temperatureSensor/api/grpc/gen"
	"google.golang.org/grpc"
)

type device interface {
	GenerateData() int32
}

type TemperatureSensorDeps struct {
	Device device
}

type TemperatureSensor struct {
	pb.UnimplementedTemperatureSensorServiceServer
	device device
}

func New(d *TemperatureSensorDeps) *TemperatureSensor {
	return &TemperatureSensor{
		device: d.Device,
	}
}

func (t *TemperatureSensor) Register(server *grpc.Server) {
	pb.RegisterTemperatureSensorServiceServer(server, t)
}

func (t *TemperatureSensor) GetData(ctx context.Context, request *pb.GetDataRequest) (*pb.GetDataResponse, error) {
	return &pb.GetDataResponse{
		Data: t.device.GenerateData(),
	}, nil
}
