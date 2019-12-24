package inmem

import (
	"sync"

	"gokit-ddd-demo/user_svc/domain/common"
	"gokit-ddd-demo/user_svc/domain/user"
)

var _ user.Repo = (*userRepo)(nil)

type userRepo struct {
	rw      sync.RWMutex
	users   map[int64]*user.User
	idIndex int64
}

func NewUserRepo() *userRepo {
	return &userRepo{
		users:   make(map[int64]*user.User),
		idIndex: 1,
	}
}

func (r *userRepo) InitMockData(users []*user.User) error {
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

func (r *userRepo) Find() ([]*user.User, error) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	items := make([]*user.User, 0, len(r.users))
	for _, u := range r.users {
		uclone := *u
		items = append(items, &uclone)
	}
	return items, nil
}

func (r *userRepo) Get(id int64) (*user.User, error) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	u, ok := r.users[id]
	if !ok {
		return nil, common.ErrNotFound
	}
	uclone := *u
	return &uclone, nil
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

func (r *userRepo) Save(u *user.User) error {
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
