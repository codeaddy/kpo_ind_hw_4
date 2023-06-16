package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"hw4/internal/app/pkg/config"
)

func NewDB(ctx context.Context) (*Database, error) {
	//dsn := generateDsn()
	dsn := "postgresql://test:test@postgres/postgres?sslmode=disable"
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return NewDatabase(pool), nil
}

func generateDsn() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Dbname)
}
