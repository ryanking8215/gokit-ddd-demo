package api

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"gokit-ddd-demo/lib"
	"gokit-ddd-demo/order_svc/svc/order"
)

type Response struct {
	Value interface{}
	Error error
}

func MakeFindEndpoint(s order.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		userID, ok := request.(int64)
		if !ok {
			return Response{nil, lib.NewError(lib.ErrInvalidArgument, "invalid user id")}, nil
		}
		v, err := s.Find(ctx, userID)
		return Response{v, err}, nil
	}
}

func MakeGetEndpoint(s order.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		id := request.(int64)
		v, err := s.Get(ctx, id)
		return v, err
	}
}
