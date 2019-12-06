package user

import (
	"context"

	"gokit-ddd-demo/user_svc/domain/models"
)

type Service interface {
	Find(ctx context.Context) ([]*models.User, error)
	Get(ctx context.Context, id int64) (*models.User, error)
	Delete(ctx context.Context, id int64) error
	Save(ctx context.Context, u *models.User) error
}

// compile time assertion shows 'userService' implements interface 'Service'
var _ Service = (*userService)(nil)

type userService struct {
	repo Repo
}

func NewUserService(repo Repo) *userService {
	return &userService{repo: repo}
}

func (s *userService) Find(ctx context.Context) ([]*models.User, error) {
	users, err := s.repo.Find()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *userService) Get(ctx context.Context, id int64) (*models.User, error) {
	u, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *userService) Delete(ctx context.Context, id int64) error {
	err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) Save(ctx context.Context, u *models.User) error {
	return s.repo.Save(u)
}
