package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"hw4/internal/app/pkg/db"
	"hw4/internal/app/pkg/user"
	"time"
)

var (
	ErrObjectNotFound = errors.New("object not found")
)

type UserRepo struct {
	db db.DBops
}

func NewUser(db db.DBops) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user user.User) (int, error) {
	var id int
	err := r.db.ExecQueryRow(ctx, "INSERT INTO public.user (username, email, password_hash, role) VALUES ($1, $2, $3, $4) RETURNING id",
		user.Username, user.Email, user.PasswordHash, user.Role).Scan(&id)
	fmt.Println(err)
	return id, err
}

func (r *UserRepo) GetById(ctx context.Context, id int) (user.User, error) {
	var u user.User
	err := r.db.Get(ctx, &u, "SELECT id,username,email,password_hash,role,created_at,updated_at FROM public.user WHERE id=$1", id)
	if err == sql.ErrNoRows {
		return user.User{}, ErrObjectNotFound
	}
	return u, err
}

func (r *UserRepo) GetByEmailPassword(ctx context.Context, email string, passwordHash string) (user.User, error) {
	var u user.User
	err := r.db.Get(ctx, &u, "SELECT id,username,email,password_hash,role,created_at,updated_at FROM public.user WHERE email=$1", email)
	if err == sql.ErrNoRows {
		return user.User{}, ErrObjectNotFound
	}
	return u, err
}

func (r *UserRepo) Update(ctx context.Context, user user.User) (bool, error) {
	result, err := r.db.Exec(ctx,
		"UPDATE public.user SET username = $1, email = $2, password_hash = $3, role = $4, updated_at = $5 WHERE id = $6",
		user.Username, user.Email, user.PasswordHash, user.Role, time.Now(), user.ID)
	return result.RowsAffected() > 0, err
}

func (r *UserRepo) Delete(ctx context.Context, id int) (bool, error) {
	result, err := r.db.Exec(ctx,
		"DELETE FROM public.user WHERE id = $1", id)
	return result.RowsAffected() > 0, err
}
