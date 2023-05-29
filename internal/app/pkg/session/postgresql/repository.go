package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"hw4/internal/app/pkg/db"
	"hw4/internal/app/pkg/session"
)

var (
	ErrObjectNotFound = errors.New("object not found")
)

type SessionRepo struct {
	db db.DBops
}

func NewSession(db db.DBops) *SessionRepo {
	return &SessionRepo{db: db}
}

func (r *SessionRepo) Create(ctx context.Context, session session.Session) (int, error) {
	var id int
	err := r.db.ExecQueryRow(ctx, "INSERT INTO public.session(user_id, session_token, expires_at) VALUES ($1, $2, $3) RETURNING ID",
		session.UserID, session.SessionToken, session.ExpiresAt).Scan(&id)
	return id, err
}

func (r *SessionRepo) GetById(ctx context.Context, id int) (session.Session, error) {
	var s session.Session
	err := r.db.Get(ctx, &s, "SELECT id,user_id,session_token,expires_at FROM public.session WHERE id=$1", id)
	if err == sql.ErrNoRows {
		return session.Session{}, ErrObjectNotFound
	}
	return s, err
}

func (r *SessionRepo) GetByUserId(ctx context.Context, userID int) (session.Session, error) {
	var s session.Session
	err := r.db.Get(ctx, &s, "SELECT id,user_id,session_token,expires_at FROM public.session WHERE user_id=$1", userID)
	if err == sql.ErrNoRows {
		return session.Session{}, ErrObjectNotFound
	}
	return s, err
}

func (r *SessionRepo) GetByToken(ctx context.Context, token string) (session.Session, error) {
	var s session.Session
	err := r.db.Get(ctx, &s, "SELECT id,user_id,session_token,expires_at FROM public.session WHERE session_token=$1", token)
	if err == sql.ErrNoRows {
		return session.Session{}, ErrObjectNotFound
	}
	return s, err
}

func (r *SessionRepo) Update(ctx context.Context, session session.Session) (bool, error) {
	result, err := r.db.Exec(ctx,
		"UPDATE public.session SET user_id = $1, session_token = $2, expires_at = $3 WHERE id = $4",
		session.UserID, session.SessionToken, session.ExpiresAt, session.ID)
	return result.RowsAffected() > 0, err
}

func (r *SessionRepo) Delete(ctx context.Context, id int) (bool, error) {
	result, err := r.db.Exec(ctx,
		"DELETE FROM public.session WHERE id = $1", id)
	return result.RowsAffected() > 0, err
}
