package coffeezone

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	openOnce sync.Once
	connStr  string
	ctx      context.Context
	pool     *pgxpool.Pool
	err      error
}

var database *Database = &Database{}

func checkConnection() error {
	bg := context.Background()
	conn, err := pgx.Connect(bg, database.connStr)
	if err != nil {
		return err
	}

	conn.Close(bg)
	return nil
}

func SetAndCheckConn(str string) error {
	database.connStr = str
	return checkConnection()
}

func DatabasePoolAvailable() bool {
	return database.err == nil && database.pool != nil
}

func NewDatabasePool() (*pgxpool.Pool, error) {
	database.openOnce.Do(func() {
		database.ctx = context.Background()
		database.pool, database.err = pgxpool.New(database.ctx, database.connStr)
	})

	if database.err != nil {
		return nil, database.err
	}

	return database.pool, nil
}

type QueryExecFunc func(pgx.Rows, error)

func QueryExec(fn QueryExecFunc, sql string, args ...any) {
	pool, err := NewDatabasePool()
	if err != nil {
		return
	}

	fn(pool.Query(database.ctx, sql, args...))
}

type QueryRowExecFunc func(pgx.Row)

func QueryRowExec(fn QueryRowExecFunc, sql string, args ...any) {
	pool, err := NewDatabasePool()
	if err != nil {
		return
	}

	fn(pool.QueryRow(database.ctx, sql, args...))
}

func CloseDatabasePool() {
	if DatabasePoolAvailable() {
		database.pool.Close()
	}
}

type Sql interface {
	CreateOrUpdate() (string, []any)
}
