package coffeezone

import (
	"context"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	openOnce  sync.Once
	connStr   string
	ctx       context.Context
	ctxCancel context.CancelFunc
	pool      *pgxpool.Pool
	err       error
}

var database *Database = &Database{}

func SetConn(str string) {
	database.connStr = str
}

func DatabasePoolAvailable() bool {
	return database.err == nil && database.pool != nil
}

func NewDatabasePool() (*pgxpool.Pool, error) {
	database.openOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		database.ctx = ctx
		database.ctxCancel = cancel

		pool, err := pgxpool.New(ctx, database.connStr)
		database.pool = pool
		database.err = err
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
		database.ctxCancel()
		database.pool.Close()
	}
}

type Sql interface {
	CreateOrUpdate() (string, []any)
}
