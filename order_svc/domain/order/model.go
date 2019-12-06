package order

type Order struct {
	ID      int64  `json:"id"`
	UserID  int64  `json:"user_id"`
	Product string `json:"product"`
}
