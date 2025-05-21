package helper

import (
	"context"
	"database/sql"
	"fmt"
	"runtime/debug"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/mattn/go-sqlite3"
)

type (
	DB = map[string]*sql.DB
	TX = map[string]*sql.Tx
)

// DBConfig 数据库初始化配置
type DBConfig struct {
	// Driver 驱动名称
	Driver string
	// DSN 数据源名称
	//
	//  [-- MySQL] username:password@tcp(localhost:3306)/dbname?timeout=10s&charset=utf8mb4&collation=utf8mb4_general_ci&parseTime=True&loc=Local
	//  [Postgres] host=localhost port=5432 user=root password=secret dbname=test search_path=schema connect_timeout=10 sslmode=disable
	//  [- SQLite] file::memory:?cache=shared
	DSN string
	// MaxOpenConns 设置最大可打开的连接数
	MaxOpenConns int
	// MaxIdleConns 连接池最大闲置连接数
	MaxIdleConns int
	// ConnMaxLifetime 连接的最大生命时长
	ConnMaxLifetime time.Duration
	// ConnMaxIdleTime 连接最大闲置时间
	ConnMaxIdleTime time.Duration
}

// NewDB returns a new sql.DB
func NewDB(cfg *DBConfig) (*sql.DB, error) {
	db, err := sql.Open(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	return db, nil
}

// Transaction 执行数据库事务
func Transaction(ctx context.Context, db *sql.DB, fn func(ctx context.Context, tx *sql.Tx) error, opts ...*sql.TxOptions) (err error) {
	var opt *sql.TxOptions
	if len(opts) != 0 {
		opt = opts[0]
	}

	tx, _err := db.BeginTx(ctx, opt)
	if _err != nil {
		err = fmt.Errorf("begin transaction: %w", _err)
		return
	}

	rollback := func(err error) error {
		if e := tx.Rollback(); e != nil {
			err = fmt.Errorf("%w: rollback: %w", err, e)
		}
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			// if panic, should rollback
			e := fmt.Errorf("transaction panic recovered: %+v", r)
			err = fmt.Errorf("%w\n%s", rollback(e), string(debug.Stack()))
		}
	}()

	if e := fn(ctx, tx); e != nil {
		err = rollback(e)
		return
	}

	if e := tx.Commit(); e != nil {
		err = fmt.Errorf("commit: %w", e)
	}
	return
}

// TransactionX 执行多数据库事务
func TransactionX(ctx context.Context, db DB, fn func(ctx context.Context, tx TX) error, opts ...*sql.TxOptions) (err error) {
	var opt *sql.TxOptions
	if len(opts) != 0 {
		opt = opts[0]
	}

	tx := make(TX, len(db))
	for k, v := range db {
		x, e := v.BeginTx(ctx, opt)
		if e != nil {
			err = fmt.Errorf("begin transaction (%s): %w", k, e)
			return
		}
		tx[k] = x
	}

	rollback := func(err error) error {
		for k, v := range tx {
			if e := v.Rollback(); e != nil {
				err = fmt.Errorf("%w; %s.Rollback: %w", err, k, e)
			}
		}
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			// if panic, should rollback
			e := fmt.Errorf("transaction panic recovered: %+v", r)
			err = fmt.Errorf("%w\n%s", rollback(e), string(debug.Stack()))
		}
	}()

	if e := fn(ctx, tx); e != nil {
		err = rollback(e)
		return
	}

	for k, v := range tx {
		if e := v.Commit(); e != nil {
			err = rollback(fmt.Errorf("%s.Commit: %w", k, e))
			return
		}
	}
	return
}
