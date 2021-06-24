package api

import (
	"context"
	"gokit-ddd-demo/api_gateway/svc"
	"gokit-ddd-demo/order_svc/svc/order"
	"gokit-ddd-demo/user_svc/svc/user"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd/lb"
)

func MakeFindUsersEndpoint(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (reponse interface{}, err error) {
		withOrders := request.(bool)
		users, err := s.UserService().Find(ctx)
		if err != nil {
			return nil, err
		}

		usersRsp := make([]*User, 0, len(users))
		for _, u := range users {
			usersRsp = append(usersRsp, &User{u, nil})
		}

		if withOrders {
			for _, u := range usersRsp {
				orders, err := s.OrderService().Find(ctx, u.ID)
				if err != nil {
					return nil, err
				}
				u.Orders = orders
			}
		}

		return usersRsp, nil
	}
}

func MakeGetUserEndpoint(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (reponse interface{}, err error) {
		id := request.(int64)
		user, err := s.UserService().Get(ctx, id)
		if err != nil {
			if retryErr, ok := err.(lb.RetryError); ok {
				return nil, retryErr.Final
			}
			return nil, err
		}

		return user, nil
	}
}

type User struct {
	*user.User
	Orders []*order.Order `json:"orders,omitempty"`
}
