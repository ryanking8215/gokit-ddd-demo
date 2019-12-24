package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"gokit-ddd-demo/order_svc/domain/order"
)

func MakeFindEndpoint(s order.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		userID := int64(0)
		if request != nil {
			userID = request.(int64)
		}
		v, err := s.Find(ctx, userID)
		return v, err
	}
}

func MakeGetEndpoint(s order.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		id := request.(int64)
		v, err := s.Get(ctx, id)
		return v, err
	}
}
