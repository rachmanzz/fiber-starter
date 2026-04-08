package cores

import (
	"context"
	"fmt"
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
		dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			Config().Database.User,
			Config().Database.Password,
			Config().Database.Host,
			Config().Database.Port,
			Config().Database.Name,
			Config().Database.SSLMode,
		)

		config, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			zap.L().Fatal("Failed to parse database DSN", zap.Error(err))
		}

		config.MaxConns = 10
		config.MinConns = 2
		config.MaxConnLifetime = 1 * time.Hour
		config.MaxConnIdleTime = 30 * time.Minute

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
