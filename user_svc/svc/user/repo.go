package user

type Repo interface {
	Find() ([]*User, error)
	Get(id int64) (*User, error)
}
