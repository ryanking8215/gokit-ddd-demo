package user

import (
	"gokit-ddd-demo/user_svc/domain/models"
)

type Repo interface {
	Find() ([]*models.User, error)
	Get(id int64) (*models.User, error)
	Delete(id int64) error
	Save(u *models.User) error
}
