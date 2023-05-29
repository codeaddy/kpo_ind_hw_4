package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"hw4/internal/app/pkg/db"
	"hw4/internal/app/pkg/dish"
	"time"
)

var (
	ErrObjectNotFound = errors.New("object not found")
)

type DishRepo struct {
	db db.DBops
}

func NewDish(db db.DBops) *DishRepo {
	return &DishRepo{db: db}
}

func (r *DishRepo) Create(ctx context.Context, dish dish.Dish) (int, error) {
	var id int
	err := r.db.ExecQueryRow(ctx, "INSERT INTO public.dish(name, description, price, quantity, is_available) VALUES ($1, $2, $3, $4, $5) RETURNING ID",
		dish.Name, dish.Description, dish.Price, dish.Quantity, dish.IsAvailable).Scan(&id)
	return id, err
}

func (r *DishRepo) GetById(ctx context.Context, id int) (dish.Dish, error) {
	var d dish.Dish
	err := r.db.Get(ctx, &d, "SELECT id,name,description,price,quantity,is_available,created_at,updated_at FROM public.dish WHERE id=$1", id)
	if err == sql.ErrNoRows {
		return dish.Dish{}, ErrObjectNotFound
	}
	return d, err
}

func (r *DishRepo) GetAll(ctx context.Context) ([]*dish.Dish, error) {
	d := make([]*dish.Dish, 0)
	err := r.db.Select(ctx, &d, "SELECT id,name,description,price,quantity,is_available,created_at,updated_at FROM public.dish")
	return d, err
}

func (r *DishRepo) Update(ctx context.Context, dish dish.Dish) (bool, error) {
	result, err := r.db.Exec(ctx,
		"UPDATE public.dish SET name = $1, description = $2, price = $3, quantity = $4, is_available = $5, updated_at = $6 WHERE id = $7",
		dish.Name, dish.Description, dish.Price, dish.Quantity, dish.IsAvailable, time.Now(), dish.ID)
	return result.RowsAffected() > 0, err
}

func (r *DishRepo) Delete(ctx context.Context, id int) (bool, error) {
	result, err := r.db.Exec(ctx,
		"DELETE FROM public.dish WHERE id = $1", id)
	return result.RowsAffected() > 0, err
}
