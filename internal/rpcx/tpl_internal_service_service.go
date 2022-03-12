package rpcx

var internalServiceService = `
package service

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

	"{{ .Package }}/api/gen/go/api"
	"github.com/hatlonely/go-kit/rpcx"
)

type Options struct {
	SleepTime time.Duration
}

func New{{ .Service }}ServiceWithOptions(options *Options) (*{{ .Service }}Service, error) {
	return &{{ .Service }}Service{
		options: options,
	}, nil
}

type {{ .Service }}Service struct {
	api.Unsafe{{ .Service }}ServiceServer

	options *Options
}

func (s *{{ .Service }}Service) Echo(ctx context.Context, req *api.EchoReq) (*api.EchoRes, error) {
	time.Sleep(s.options.SleepTime)

	header := metadata.Pairs("Location", "https://www.baidu.com")
	if err := grpc.SendHeader(ctx, header); err != nil {
		return nil, errors.Wrap(err, "grpc.SendHeader failed")
	}

	return &api.EchoRes{
		Message: req.Message,
	}, nil
}

func (s *{{ .Service }}Service) Add(ctx context.Context, req *api.AddReq) (*api.AddRes, error) {
	time.Sleep(s.options.SleepTime)

	if req.I1 > 100 || req.I2 > 100 {
		return nil, rpcx.NewError(errors.Errorf("parameter too large"), codes.InvalidArgument, "InvalidArgument", "parameter too large")
	}

	return &api.AddRes{
		Val: req.I1 + req.I2,
	}, nil
}
`
