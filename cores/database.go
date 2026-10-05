package cores

import (
	"context"
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
	dbOnce       sync.Once
	contractFn   func(*pgxpool.Pool)
	contractOnce sync.Once
)

func ConnectDB() {
	dbOnce.Do(func() {
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
			zap.L().Fatal("Failed to parse database DSN", zap.Error(err))
		}

		config.MaxConns = Config().Database.MaxConns
		config.MinConns = Config().Database.MinConns
		config.MaxConnLifetime = Config().Database.MaxConnLifetime
		config.MaxConnIdleTime = Config().Database.MaxConnIdleTime

		pool, err := pgxpool.NewWithConfig(context.Background(), config)
		if err != nil {
			zap.L().Fatal("Failed to connect to database", zap.Error(err))
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			zap.L().Fatal("Database ping failed", zap.Error(err))
		}

		zap.L().Info("Database connection established")

		db = pool

		if contractFn != nil {
			contractFn(pool)
		}
	})
}

func SetDatabaseContract(fn func(*pgxpool.Pool)) {
	contractOnce.Do(func() {
		contractFn = fn
	})
}

func CloseDB() {
	if db != nil {
		db.Close()
		zap.L().Info("Database connection pool closed")
	}
}
