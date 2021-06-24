package inmem

import (
	"sync"

	"gokit-ddd-demo/lib"
	"gokit-ddd-demo/order_svc/svc/order"
)

var _ order.Repo = (*orderRepo)(nil)

type orderRepo struct {
	rw     sync.RWMutex
	orders map[int64]*order.Order
	// userIndex map[int64]int64 // userID ->
	idIndex int64
}

func NewOrderRepo() *orderRepo {
	return &orderRepo{
		orders:  make(map[int64]*order.Order),
		idIndex: 1,
	}
}

func (r *orderRepo) InitMockData(orders []*order.Order) error {
	r.rw.Lock()
	defer r.rw.Unlock()

	for _, o := range orders {
		r.orders[o.ID] = o
		if o.ID > r.idIndex {
			r.idIndex = o.ID
		}
	}
	return nil
}

func (r *orderRepo) Find(userID int64) ([]*order.Order, error) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	items := make([]*order.Order, 0, len(r.orders))
	for _, o := range r.orders {
		if userID <= 0 {
			oclone := *o
			items = append(items, &oclone)
		} else {
			if o.UserID == userID {
				oclone := *o
				items = append(items, &oclone)
			}
		}
	}
	return items, nil
}

func (r *orderRepo) Get(id int64) (*order.Order, error) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	o, ok := r.orders[id]
	if !ok {
		return nil, lib.NewError(lib.ErrNotFound, "")
	}
	oclone := *o
	return &oclone, nil
}

func (r *orderRepo) Delete(id int64) error {
	r.rw.Lock()
	defer r.rw.Unlock()

	_, ok := r.orders[id]
	if !ok {
		return lib.NewError(lib.ErrNotFound, "")
	}
	delete(r.orders, id)

	return nil
}

func (r *orderRepo) Save(o *order.Order) error {
	r.rw.Lock()
	defer r.rw.Unlock()

	_, ok := r.orders[o.ID]
	if !ok {
		o.ID = r.idIndex
		r.idIndex++
	}
	r.orders[o.ID] = o

	return nil
}
