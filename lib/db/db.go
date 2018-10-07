package db

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// TODO Transactions??

type DBHandler interface {
	DB() *sqlx.DB
	Host() string
	Name() string
	Port() string

	Ping() error
	Close()

	GetById(ctx context.Context, result interface{}, query string, id interface{}) error
	Select(ctx context.Context, results interface{}, query string, params ...interface{}) error
	InsertWithId(ctx context.Context, query string, params ...interface{}) (int64, error)
	Insert(ctx context.Context, query string, params ...interface{}) error
	Update(ctx context.Context, query string, params ...interface{}) (int64, error)
	Delete(ctx context.Context, query string, params ...interface{}) (int64, error)
}
