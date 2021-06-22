package svc

import (
	"gokit-ddd-demo/order_svc/svc/order"
	"gokit-ddd-demo/user_svc/svc/user"
)

type Service interface {
	UserService() user.Service
	OrderService() order.Service
}

var _ Service = (*service)(nil)

type service struct {
	userSvc  user.Service
	orderSvc order.Service
}

func (s *service) UserService() user.Service {
	return s.userSvc
}

func (s *service) OrderService() order.Service {
	return s.orderSvc
}

func NewService(userSvc user.Service, orderSvc order.Service) Service {
	return &service{userSvc, orderSvc}
}
