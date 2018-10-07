package psql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	SSLModeDisable SSLMode = "disable"     // No SSL
	SSLModeRequire SSLMode = "require"     // Always SSL (skip verification)
	SSLModeFull    SSLMode = "verify-full" // Always SSL (require verification)
)

type SSLMode string

type PSQLHandler struct {
	conn *sqlx.DB
	host string
	name string
	port int
}

// Postgres Connection Object.  Wraps a sqlx.DB and provides a DBConnection() method to access the
// DB Connection
func NewPSQLDBHandler(host, dbname, user, password string, port int, sslmode SSLMode) (*PSQLHandler, error) {
	dsn := fmt.Sprintf("host=%s user=%s password='%s' dbname=%s port=%d sslmode=%s",
		host,
		user,
		password,
		dbname,
		port,
		sslmode,
	)
	conn, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	conn.SetMaxOpenConns(20)

	psql := &PSQLHandler{
		conn: conn,
		host: host,
		name: dbname,
		port: port,
	}

	err = psql.Ping()
	if err != nil {
		return nil, err
	}

	return psql, nil
}

func (t *PSQLHandler) DB() *sqlx.DB {
	return t.conn
}

func (t *PSQLHandler) Host() string {
	return t.host
}

func (t *PSQLHandler) Name() string {
	return t.name
}

func (t *PSQLHandler) Port() string {
	return fmt.Sprintf("%d", t.port)
}

func (t *PSQLHandler) Ping() error {
	return t.conn.Ping()
}

func (t *PSQLHandler) Close() {
	if err := t.conn.Close(); err != nil {
		panic(err)
	}
}

func IsDuplicateKey(err error) bool {
	return strings.Contains(err.Error(), "duplicate key")
}

func (t *PSQLHandler) GetById(ctx context.Context, result interface{}, query string, id interface{}) error {

	tx, err := t.conn.Beginx()
	if err != nil {
		return err
	}

	if err := tx.Get(result, query, id); err != nil {
		if sql.ErrNoRows == err {
			return Commit(ctx, tx, nil)
		}
		return Rollback(ctx, tx, err)
	}

	return Commit(ctx, tx, err)
}

func (t *PSQLHandler) Select(ctx context.Context, results interface{}, query string, params ...interface{}) error {

	tx, err := t.conn.Beginx()
	if err != nil {
		return err
	}

	if err := tx.Select(results, query, params...); err != nil {
		if sql.ErrNoRows == err {
			return Commit(ctx, tx, nil)
		}
		return Rollback(ctx, tx, err)
	}

	return Commit(ctx, tx, err)
}

func (t *PSQLHandler) InsertWithId(ctx context.Context, query string, params ...interface{}) (int64, error) {

	if !strings.Contains(strings.ToUpper(query), "RETURNING") {
		panic(fmt.Sprintf("Query (%s) needs to contain a 'RETURNING id' expression", query))
	}

	tx, err := t.conn.Beginx()
	if err != nil {
		return 0, err
	}

	var id int64
	if err := tx.QueryRow(query, params...).Scan(&id); sql.ErrNoRows == err {
		return 0, err
	}

	return id, err
}

func (t *PSQLHandler) Insert(ctx context.Context, query string, params ...interface{}) error {

	tx, err := t.conn.Beginx()
	if err != nil {
		return err
	}

	// TODO result
	if _, err := tx.Exec(query, params...); sql.ErrNoRows == err {
		return err
	}
	return err
}

func (t *PSQLHandler) Update(ctx context.Context, query string, params ...interface{}) (int64, error) {

	tx, err := t.conn.Beginx()
	if err != nil {
		return 0, err
	}

	result, err := tx.Exec(query, params...)
	if err != nil {
		return 0, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return 0, nil
	}
	return affected, nil
}

func (t *PSQLHandler) Delete(ctx context.Context, query string, params ...interface{}) (int64, error) {
	return t.Update(ctx, query, params...)
}

func Commit(ctx context.Context, tx *sqlx.Tx, err error) error {
	if err != nil {
		return Rollback(ctx, tx, err)
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func Rollback(ctx context.Context, tx *sqlx.Tx, err error) error {
	if err := tx.Rollback(); err != nil {
		return err
	}
	return err
}
