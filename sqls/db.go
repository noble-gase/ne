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
	//  [-MySQL] mysql
	//  [-PgSQL] pgx
	//  [SQLite] sqlite3
	Driver string
	// DSN 数据源名称
	//
	//  [-MySQL] username:password@tcp(host:3306)/dbname?timeout=10s&charset=utf8mb4&parseTime=True&loc=Local
	//  [-PgSQL] postgres://username:password@host:5432/dbname
	//  [SQLite] file::memory:?cache=shared || file:/path/test.db
	DSN string
	// MaxOpenConns 设置最大可打开的连接数
	MaxOpenConns int
	// MaxIdleConns 连接池最大闲置连接数
	MaxIdleConns int
	// ConnMaxIdleTime 连接最大闲置时间
	ConnMaxIdleTime time.Duration
	// ConnMaxLifetime 连接的最大生命时长
	ConnMaxLifetime time.Duration
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
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return db, nil
}
