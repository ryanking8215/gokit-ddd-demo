package order

type Repo interface {
	Find(userID int64) ([]*Order, error)
	Get(id int64) (*Order, error)
	Delete(id int64) error
	Save(o *Order) error
}
