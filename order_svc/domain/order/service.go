package order

import (
	"context"
)

type Service interface {
	Find(ctx context.Context, userID int64) ([]*Order, error)
	Get(ctx context.Context, id int64) (*Order, error)
	Delete(ctx context.Context, id int64) error
	Save(ctx context.Context, o *Order) error
}

// compile time assertion shows 'userService' implements interface 'Service'
var _ Service = (*orderService)(nil)

type orderService struct {
	repo Repo
}

func NewOrderService(repo Repo) *orderService {
	return &orderService{repo: repo}
}

func (s *orderService) Find(ctx context.Context, userID int64) ([]*Order, error) {
	orders, err := s.repo.Find(userID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *orderService) Get(ctx context.Context, id int64) (*Order, error) {
	o, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (s *orderService) Delete(ctx context.Context, id int64) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *orderService) Save(ctx context.Context, o *Order) error {
	return s.repo.Save(o)
}
