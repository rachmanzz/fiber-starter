package cores

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rachmanzz/fiber-starter/config"
	"go.uber.org/zap"
)

// DefaultDatabaseName is the default identifier for the primary PostgreSQL connection pool.
const DefaultDatabaseName = "default"

var (
	dbMu            sync.RWMutex
	db              *pgxpool.Pool
	pools           = make(map[string]*pgxpool.Pool)
	contractFn      func(*pgxpool.Pool)
	namedContractFn func(name string, pool *pgxpool.Pool)
)

// ConnectDB initializes the default PostgreSQL connection pool and returns an error on failure.
func ConnectDB() error {
	_, err := ConnectNamedDB(DefaultDatabaseName, Config().Database)
	return err
}

// ConnectNamedDB initializes a named PostgreSQL connection pool with the provided configuration.
func ConnectNamedDB(name string, cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	dbMu.Lock()
	if pool, exists := pools[name]; exists && pool != nil {
		dbMu.Unlock()
		return pool, nil
	}
	dbMu.Unlock()

	dsn := cfg.DSN()
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database DSN for %q: %w", name, err)
	}

	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = cfg.MinConns
	poolConfig.MaxConnLifetime = cfg.MaxConnLifetime
	poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create database connection pool for %q: %w", name, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("database ping failed for %q: %w", name, err)
	}

	dbMu.Lock()
	// Check again in case of race condition during async connection initialization
	if existing, exists := pools[name]; exists && existing != nil {
		dbMu.Unlock()
		pool.Close()
		return existing, nil
	}

	pools[name] = pool
	if name == DefaultDatabaseName {
		db = pool
	}
	cFn := contractFn
	ncFn := namedContractFn
	dbMu.Unlock()

	zap.L().Info("Database connection established", zap.String("pool", name))

	if name == DefaultDatabaseName && cFn != nil {
		cFn(pool)
	}
	if ncFn != nil {
		ncFn(name, pool)
	}

	return pool, nil
}

// RegisterDBPool registers an existing pgxpool.Pool under a designated pointer name.
func RegisterDBPool(name string, pool *pgxpool.Pool) {
	dbMu.Lock()
	pools[name] = pool
	if name == DefaultDatabaseName {
		db = pool
	}
	cFn := contractFn
	ncFn := namedContractFn
	dbMu.Unlock()

	if name == DefaultDatabaseName && cFn != nil {
		cFn(pool)
	}
	if ncFn != nil {
		ncFn(name, pool)
	}
}

// HasDB checks if a named database pool is registered.
func HasDB(name string) bool {
	dbMu.RLock()
	defer dbMu.RUnlock()
	p, exists := pools[name]
	return exists && p != nil
}


// SetDatabaseContract registers the database consumer contract callback for the default pool.
// If the default pool is already connected, the callback is invoked immediately.
func SetDatabaseContract(fn func(*pgxpool.Pool)) {
	dbMu.Lock()
	contractFn = fn
	currentDB := db
	dbMu.Unlock()

	if currentDB != nil && fn != nil {
		fn(currentDB)
	}
}

// SetNamedDatabaseContract registers a callback triggered whenever any named database pool connects.
// If pools are already registered, the callback is invoked for each existing pool immediately.
func SetNamedDatabaseContract(fn func(name string, pool *pgxpool.Pool)) {
	dbMu.Lock()
	namedContractFn = fn
	existing := make(map[string]*pgxpool.Pool, len(pools))
	for k, v := range pools {
		existing[k] = v
	}
	dbMu.Unlock()

	if fn != nil {
		for k, v := range existing {
			fn(k, v)
		}
	}
}

// CloseDB closes all active PostgreSQL connection pools safely.
func CloseDB() {
	dbMu.Lock()
	toClose := make(map[string]*pgxpool.Pool, len(pools))
	for k, v := range pools {
		toClose[k] = v
	}
	pools = make(map[string]*pgxpool.Pool)
	db = nil
	contractFn = nil
	namedContractFn = nil
	dbMu.Unlock()

	for name, pool := range toClose {
		if pool != nil {
			func() {
				defer func() { _ = recover() }()
				pool.Close()
			}()
			zap.L().Info("Database connection pool closed", zap.String("pool", name))
		}
	}
}

// CloseNamedDB closes a specific named PostgreSQL connection pool.
func CloseNamedDB(name string) error {
	dbMu.Lock()
	pool, exists := pools[name]
	if !exists || pool == nil {
		dbMu.Unlock()
		return fmt.Errorf("database pool %q does not exist", name)
	}
	delete(pools, name)
	if name == DefaultDatabaseName {
		db = nil
	}
	dbMu.Unlock()

	func() {
		defer func() { _ = recover() }()
		pool.Close()
	}()
	zap.L().Info("Named database connection pool closed", zap.String("pool", name))
	return nil
}
