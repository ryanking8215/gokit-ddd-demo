package api

import (
	"gokit-ddd-demo/order_svc/domain/order"
	"gokit-ddd-demo/user_svc/domain/user"
)

type Order = order.Order

type User struct {
	*user.User
	Orders []*Order `json:"orders"`
}
