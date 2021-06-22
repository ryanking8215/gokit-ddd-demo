package api

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"gokit-ddd-demo/user_svc/svc/user"
)

func MakeFindEndpoint(s user.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		v, err := s.Find(ctx)
		return v, err
	}
}

func MakeGetEndpoint(s user.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		id := request.(int64)
		v, err := s.Get(ctx, id)
		return v, err
	}
}
