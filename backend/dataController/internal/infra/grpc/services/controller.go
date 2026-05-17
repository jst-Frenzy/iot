package services

import (
	"context"
	dcs "github.com/jst-Frenzy/iot/backend/dataController/api/grpc/gen"
	"google.golang.org/grpc"
)

type manager interface {
	OnPump(ctx context.Context) error
	OffPump(ctx context.Context) error
	OnFan(ctx context.Context) error
	OffFan(ctx context.Context) error
}

type ControllerDeps struct {
	Manager manager
}

type Controller struct {
	dcs.UnsafeDeviceServiceServer
	manager manager
}

func NewController(d *ControllerDeps) *Controller {
	return &Controller{
		manager: d.Manager,
	}
}

func (a *Controller) Register(server *grpc.Server) {
	dcs.RegisterDeviceServiceServer(server, a)
}

func (a *Controller) OnPump(ctx context.Context, _ *dcs.OnPumpRequest) (*dcs.OnPumpResponse, error) {
	return &dcs.OnPumpResponse{}, a.manager.OnPump(ctx)
}

func (a *Controller) OffPump(ctx context.Context, _ *dcs.OffPumpRequest) (*dcs.OffPumpResponse, error) {
	return &dcs.OffPumpResponse{}, a.manager.OffPump(ctx)
}

func (a *Controller) OnFan(ctx context.Context, _ *dcs.OnFanRequest) (*dcs.OnFanResponse, error) {
	return &dcs.OnFanResponse{}, a.manager.OnFan(ctx)
}

func (a *Controller) OffFan(ctx context.Context, _ *dcs.OffFanRequest) (*dcs.OffFanResponse, error) {
	return &dcs.OffFanResponse{}, a.manager.OffFan(ctx)
}
