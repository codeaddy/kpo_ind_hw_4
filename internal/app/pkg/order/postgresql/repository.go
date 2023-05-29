package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"hw4/internal/app/pkg/db"
	"hw4/internal/app/pkg/order"
	"time"
)

var (
	ErrObjectNotFound = errors.New("object not found")
)

type OrderRepo struct {
	db db.DBops
}

func NewOrder(db db.DBops) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) Create(ctx context.Context, order order.Order) (int, error) {
	var id int
	err := r.db.ExecQueryRow(ctx, "INSERT INTO public.order(user_id, status, special_requests) VALUES ($1, $2, $3) RETURNING ID",
		order.UserID, order.Status, order.SpecialRequests).Scan(&id)
	return id, err
}

func (r *OrderRepo) GetById(ctx context.Context, id int) (order.Order, error) {
	var o order.Order
	err := r.db.Get(ctx, &o, "SELECT id,user_id,status,special_requests,created_at,updated_at FROM public.order WHERE id=$1", id)
	if err == sql.ErrNoRows {
		return order.Order{}, ErrObjectNotFound
	}
	return o, err
}

func (r *OrderRepo) Update(ctx context.Context, order order.Order) (bool, error) {
	result, err := r.db.Exec(ctx,
		"UPDATE public.order SET user_id = $1, status = $2, special_requests = $3, updated_at = $4 WHERE id = $5",
		order.UserID, order.Status, order.SpecialRequests, time.Now(), order.ID)
	return result.RowsAffected() > 0, err
}

func (r *OrderRepo) Delete(ctx context.Context, id int) (bool, error) {
	result, err := r.db.Exec(ctx,
		"DELETE FROM public.order WHERE id = $1", id)
	return result.RowsAffected() > 0, err
}
