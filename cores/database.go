package cores

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

var (
	db           *pgxpool.Pool
	dbMu         sync.Mutex
	contractFn   func(*pgxpool.Pool)
	contractOnce sync.Once
)

// ConnectDB initializes the PostgreSQL connection pool and returns an error on failure.
func ConnectDB() error {
	dbMu.Lock()
	defer dbMu.Unlock()

	if db != nil {
		return nil
	}

	u := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(Config().Database.User, Config().Database.Password),
		Host:     net.JoinHostPort(Config().Database.Host, strconv.Itoa(Config().Database.Port)),
		Path:     Config().Database.Name,
		RawQuery: "sslmode=" + Config().Database.SSLMode,
	}
	dsn := u.String()

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("failed to parse database DSN: %w", err)
	}

	config.MaxConns = Config().Database.MaxConns
	config.MinConns = Config().Database.MinConns
	config.MaxConnLifetime = Config().Database.MaxConnLifetime
	config.MaxConnIdleTime = Config().Database.MaxConnIdleTime

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("failed to create database connection pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return fmt.Errorf("database ping failed: %w", err)
	}

	zap.L().Info("Database connection established")
	db = pool

	if contractFn != nil {
		contractFn(pool)
	}

	return nil
}

// SetDatabaseContract registers the database consumer contract callback.
func SetDatabaseContract(fn func(*pgxpool.Pool)) {
	contractOnce.Do(func() {
		contractFn = fn
	})
}

// CloseDB closes the active PostgreSQL connection pool safely.
func CloseDB() {
	dbMu.Lock()
	defer dbMu.Unlock()

	if db != nil {
		db.Close()
		db = nil
		zap.L().Info("Database connection pool closed")
	}
}
