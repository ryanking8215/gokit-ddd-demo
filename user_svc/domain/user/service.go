package user

import (
	"context"
)

type Service interface {
	Find(ctx context.Context) ([]*User, error)
	Get(ctx context.Context, id int64) (*User, error)
}

// compile time assertion shows 'userService' implements interface 'Service'
var _ Service = (*userService)(nil)

type userService struct {
	repo Repo
}

func NewUserService(repo Repo) *userService {
	return &userService{repo: repo}
}

func (s *userService) Find(ctx context.Context) ([]*User, error) {
	users, err := s.repo.Find()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *userService) Get(ctx context.Context, id int64) (*User, error) {
	u, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	return u, nil
}
