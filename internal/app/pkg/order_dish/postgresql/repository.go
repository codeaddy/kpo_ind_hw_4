package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"hw4/internal/app/pkg/db"
	"hw4/internal/app/pkg/order_dish"
)

var (
	ErrObjectNotFound = errors.New("object not found")
)

type OrderDishRepo struct {
	db db.DBops
}

func NewOrderDish(db db.DBops) *OrderDishRepo {
	return &OrderDishRepo{db: db}
}

func (r *OrderDishRepo) Create(ctx context.Context, orderDish order_dish.OrderDish) (int, error) {
	var id int
	err := r.db.ExecQueryRow(ctx, "INSERT INTO public.order_dish(order_id, dish_id, quantity, price) VALUES ($1, $2, $3, $4) RETURNING ID",
		orderDish.OrderID, orderDish.DishID, orderDish.Quantity, orderDish.Price).Scan(&id)
	return id, err
}

func (r *OrderDishRepo) GetById(ctx context.Context, id int) (order_dish.OrderDish, error) {
	var o order_dish.OrderDish
	err := r.db.Get(ctx, &o, "SELECT id,order_id,dish_id,quantity,price FROM public.order_dish WHERE id=$1", id)
	if err == sql.ErrNoRows {
		return order_dish.OrderDish{}, ErrObjectNotFound
	}
	return o, err
}

func (r *OrderDishRepo) Update(ctx context.Context, orderDish order_dish.OrderDish) (bool, error) {
	result, err := r.db.Exec(ctx,
		"UPDATE public.order_dish SET order_id = $1, dish_id = $2, quantity = $3, price = $4 WHERE id = $5",
		orderDish.OrderID, orderDish.DishID, orderDish.Quantity, orderDish.Price, orderDish.ID)
	return result.RowsAffected() > 0, err
}

func (r *OrderDishRepo) Delete(ctx context.Context, id int) (bool, error) {
	result, err := r.db.Exec(ctx,
		"DELETE FROM public.order_dish WHERE id = $1", id)
	return result.RowsAffected() > 0, err
}
