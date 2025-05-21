package sqls

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/mattn/go-sqlite3"
)

// Config 数据库初始化配置
type Config struct {
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
func NewDB(cfg *Config) (*sql.DB, error) {
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
