package order_dish

type OrderDish struct {
	ID       int     `db:"id" json:"id"`
	OrderID  int     `db:"order_id" json:"order_id"`
	DishID   int     `db:"dish_id" json:"dish_id"`
	Quantity int     `db:"quantity" json:"quantity"`
	Price    float64 `db:"price" json:"price"`
}
