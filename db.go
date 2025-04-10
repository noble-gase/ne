package ne

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

// NewDB sql.DB
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

// Transaction 执行数据库事物
func Transaction(ctx context.Context, db *sql.DB, fn func(ctx context.Context, tx *sql.Tx) error) (err error) {
	tx, _err := db.BeginTx(ctx, nil)
	if _err != nil {
		err = fmt.Errorf("begin transaction: %w", _err)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("transaction panic recovered: %+v", r)
			if err_ := tx.Rollback(); err_ != nil {
				err = fmt.Errorf("%w: rollback: %w", err, err_)
			}
			err = fmt.Errorf("%w\n%s", err, string(debug.Stack()))
		}
	}()

	if err = fn(ctx, tx); err != nil {
		if err_ := tx.Rollback(); err_ != nil {
			err = fmt.Errorf("%w: transaction rollback: %w", err, err_)
		}
		return
	}

	if _err = tx.Commit(); _err != nil {
		err = fmt.Errorf("transaction commit: %w", _err)
		return
	}
	return
}
