package api

import (
	"context"
	"gokit-ddd-demo/api_gateway/svc"
	"gokit-ddd-demo/order_svc/svc/order"
	"gokit-ddd-demo/user_svc/svc/user"

	"github.com/go-kit/kit/endpoint"
)

func MakeFindUsersEndpoint(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (reponse interface{}, err error) {
		withOrders := request.(bool)
		println(">>>>>>> ", withOrders)

		users, err := s.UserService().Find(ctx)
		if err != nil {
			return nil, err
		}

		usersRsp := make([]*User, 0, len(users))
		for _, u := range users {
			usersRsp = append(usersRsp, &User{u, nil})
		}

		if withOrders {
			println("yyyyyyy")
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

type User struct {
	*user.User
	Orders []*order.Order `json:"orders,omitempty"`
}
