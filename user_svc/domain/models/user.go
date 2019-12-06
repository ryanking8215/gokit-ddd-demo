package models

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`

	Orders []*Order `json:"orders,omitmepty"`
}
