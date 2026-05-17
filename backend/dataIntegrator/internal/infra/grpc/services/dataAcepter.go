package services

import (
	"context"
	dis "github.com/jst-Frenzy/iot/backend/dataIntegrator/api/grpc/gen"
	"google.golang.org/grpc"
)

type processor interface {
	Process(context.Context, int32, int32) error
}

type AccepterDeps struct {
	Processor processor
}

type Accepter struct {
	dis.UnimplementedDataIntegratorServiceServer
	processor processor
}

func NewAccepter(d *AccepterDeps) *Accepter {
	return &Accepter{
		processor: d.Processor,
	}
}

func (a *Accepter) Register(server *grpc.Server) {
	dis.RegisterDataIntegratorServiceServer(server, a)
}

func (a *Accepter) AcceptData(ctx context.Context, req *dis.AcceptDataRequest) (*dis.AcceptDataResponse, error) {
	return &dis.AcceptDataResponse{}, a.processor.Process(ctx, req.Wet, req.Temperature)
}
