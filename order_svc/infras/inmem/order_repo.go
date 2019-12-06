package inmem

import (
	"sync"

	"gokit-ddd-demo/order_svc/domain/common"
	"gokit-ddd-demo/order_svc/domain/order"
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
	for _, v := range r.orders {
		if userID <= 0 {
			items = append(items, v)
		} else {
			if v.UserID == userID {
				items = append(items, v)
			}
		}
	}
	return items, nil
}

func (r *orderRepo) Get(id int64) (*order.Order, error) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	item, ok := r.orders[id]
	if !ok {
		return nil, common.ErrNotFound
	}
	return item, nil
}

func (r *orderRepo) Delete(id int64) error {
	r.rw.Lock()
	defer r.rw.Unlock()

	_, ok := r.orders[id]
	if !ok {
		return common.ErrNotFound
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
