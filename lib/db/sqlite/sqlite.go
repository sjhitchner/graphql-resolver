package sqlite

import (
	"context"
	"database/sql"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteHandler struct {
	conn     *sqlx.DB
	filepath string
}

// Postgres Connection Object.  Wraps a sqlx.DB and provides a DBConnection() method to access the
// DB Connection
func NewSQLiteDBHandler(filepath string) (*SQLiteHandler, error) {
	conn, err := sqlx.Connect("sqlite3", filepath)
	if err != nil {
		return nil, err
	}
	conn.SetMaxOpenConns(20)

	s := &SQLiteHandler{
		conn:     conn,
		filepath: filepath,
	}

	err = s.Ping()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (t *SQLiteHandler) DB() *sqlx.DB {
	return t.conn
}

func (t *SQLiteHandler) Host() string {
	return t.filepath
}

func (t *SQLiteHandler) Name() string {
	return "sqlite"
}

func (t *SQLiteHandler) Port() string {
	return "0"
}

func (t *SQLiteHandler) Ping() error {
	return t.conn.Ping()
}

func (t *SQLiteHandler) Close() {
	if err := t.conn.Close(); err != nil {
		panic(err)
	}
}

func IsDuplicateKey(err error) bool {
	return strings.Contains(err.Error(), "duplicate key")
}

func (t *SQLiteHandler) GetById(ctx context.Context, result interface{}, query string, id interface{}) error {

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

func (t *SQLiteHandler) Select(ctx context.Context, results interface{}, query string, params ...interface{}) error {

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

func (t *SQLiteHandler) InsertWithId(ctx context.Context, query string, params ...interface{}) (int64, error) {

	/*
		if !strings.Contains(strings.ToUpper(query), "RETURNING") {
			panic(fmt.Sprintf("Query (%s) needs to contain a 'RETURNING id' expression", query))
		}
	*/

	tx, err := t.conn.Beginx()
	if err != nil {
		return 0, Rollback(ctx, tx, err)
	}

	result, err := tx.Exec(query, params...)
	if sql.ErrNoRows == err {
		return 0, Rollback(ctx, tx, err)
	}

	if err != nil {
		return 0, Rollback(ctx, tx, err)
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, Rollback(ctx, tx, err)
	}

	return lastId, Commit(ctx, tx, err)
}

func (t *SQLiteHandler) Insert(ctx context.Context, query string, params ...interface{}) error {

	tx, err := t.conn.Beginx()
	if err != nil {
		return Rollback(ctx, tx, err)
	}

	// TODO result
	if _, err := tx.Exec(query, params...); sql.ErrNoRows == err {
		return Rollback(ctx, tx, err)
	}

	return Commit(ctx, tx, err)
}

func (t *SQLiteHandler) Update(ctx context.Context, query string, params ...interface{}) (int64, error) {

	tx, err := t.conn.Beginx()
	if err != nil {
		return 0, Rollback(ctx, tx, err)
	}

	result, err := tx.Exec(query, params...)
	if err != nil {
		return 0, Rollback(ctx, tx, err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return 0, Rollback(ctx, tx, err)
	}

	return affected, Commit(ctx, tx, err)
}

func (t *SQLiteHandler) Delete(ctx context.Context, query string, params ...interface{}) (int64, error) {
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
