package api

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"gokit-ddd-demo/user_svc/svc/user"
)

type Response struct {
	Value interface{}
	Error error
}

func MakeFindEndpoint(s user.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		v, err := s.Find(ctx)
		return &Response{v, err}, nil
	}
}

func MakeGetEndpoint(s user.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		id := request.(int64)
		v, err := s.Get(ctx, id)
		return &Response{v, err}, nil
	}
}
