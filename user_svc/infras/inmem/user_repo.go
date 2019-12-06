package inmem

import (
	"sync"

	"gokit-ddd-demo/user_svc/domain/common"
	"gokit-ddd-demo/user_svc/domain/models"
	"gokit-ddd-demo/user_svc/domain/user"
)

var _ user.Repo = (*userRepo)(nil)

type userRepo struct {
	rw      sync.RWMutex
	users   map[int64]*models.User
	idIndex int64
}

func NewUserRepo() *userRepo {
	return &userRepo{
		users:   make(map[int64]*models.User),
		idIndex: 1,
	}
}

func (r *userRepo) InitMockData(users []*models.User) error {
	r.rw.Lock()
	defer r.rw.Unlock()

	for _, u := range users {
		r.users[u.ID] = u
		if u.ID > r.idIndex {
			r.idIndex = u.ID
		}
	}
	return nil
}

func (r *userRepo) Find() ([]*models.User, error) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	items := make([]*models.User, 0, len(r.users))
	for _, v := range r.users {
		items = append(items, v)
	}
	return items, nil
}

func (r *userRepo) Get(id int64) (*models.User, error) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	item, ok := r.users[id]
	if !ok {
		return nil, common.ErrNotFound
	}
	return item, nil
}

func (r *userRepo) Delete(id int64) error {
	r.rw.Lock()
	defer r.rw.Unlock()

	_, ok := r.users[id]
	if !ok {
		return common.ErrNotFound
	}
	delete(r.users, id)

	return nil
}

func (r *userRepo) Save(u *models.User) error {
	r.rw.Lock()
	defer r.rw.Unlock()

	_, ok := r.users[u.ID]
	if !ok {
		u.ID = r.idIndex
		r.idIndex++
	}
	r.users[u.ID] = u

	return nil
}
