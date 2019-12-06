package models

type Order struct {
	ID          int64  `json:"id"`
	ProductName string `json:"product_name"`
}
