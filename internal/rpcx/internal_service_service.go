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

	"github.com/hatlonely/go-kit/examples/rpcx/api/gen/go/api"
	"github.com/hatlonely/go-kit/rpcx"
)

type Options struct {
	SleepTime time.Duration
}

func NewExampleServiceWithOptions(options *Options) (*ExampleService, error) {
	return &ExampleService{
		options: options,
	}, nil
}

type ExampleService struct {
	api.ExampleServiceServer

	options *Options
}
`
