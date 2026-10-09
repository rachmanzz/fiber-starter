package contract

import (
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rachmanzz/fiber-starter/app/repository"
)

// DefaultDatabasePointer is the identifier for the default database connection.
const DefaultDatabasePointer = "default"

var (
	regMu      sync.RWMutex
	pool       *pgxpool.Pool
	queries    *repository.Queries
	poolsMap   = make(map[string]*pgxpool.Pool)
	queriesMap = make(map[string]*repository.Queries)
)

// DatabaseContract registers the primary PostgreSQL connection pool under the "default" pointer.
func DatabaseContract(db *pgxpool.Pool) {
	RegisterNamedDatabase(DefaultDatabasePointer, db)
}

// RegisterNamedDatabase registers a database connection pool and instantiates its SQLC Querier.
func RegisterNamedDatabase(name string, db *pgxpool.Pool) {
	regMu.Lock()
	defer regMu.Unlock()

	var q *repository.Queries
	if db != nil {
		q = repository.New(db)
	}
	poolsMap[name] = db
	queriesMap[name] = q

	if name == DefaultDatabasePointer {
		pool = db
		queries = q
	}
}

// RegisterNamedQueries registers a pre-instantiated *repository.Queries under a named pointer (useful for testing and mocks).
func RegisterNamedQueries(name string, q *repository.Queries) {
	regMu.Lock()
	defer regMu.Unlock()

	queriesMap[name] = q
	if name == DefaultDatabasePointer {
		queries = q
	}
}

// SetQueries overrides the default queries instance (useful for testing).
func SetQueries(q *repository.Queries) {
	RegisterNamedQueries(DefaultDatabasePointer, q)
}

// GetQueries returns the default SQLC Queries instance. Panics if not initialized.
func GetQueries() *repository.Queries {
	regMu.RLock()
	defer regMu.RUnlock()

	if queries == nil {
		panic("Repository Error: database contract must be initialized before accessing GetQueries")
	}
	return queries
}

// GetQueriesNamed retrieves the SQLC Queries instance for a specific named pointer.
func GetQueriesNamed(name string) (*repository.Queries, error) {
	regMu.RLock()
	defer regMu.RUnlock()

	target := name
	if target == "" {
		target = DefaultDatabasePointer
	}

	q, exists := queriesMap[target]
	if !exists || q == nil {
		return nil, fmt.Errorf("repository contract pointer %q not registered", target)
	}
	return q, nil
}

// Use retrieves the Queries instance for a named database pointer.
// If name is empty or "default", it returns the default queries.
// If the pointer name does not exist, it panics with a descriptive error.
func Use(name ...string) *repository.Queries {
	target := DefaultDatabasePointer
	if len(name) > 0 && name[0] != "" {
		target = name[0]
	}

	if target == DefaultDatabasePointer {
		return GetQueries()
	}

	q, err := GetQueriesNamed(target)
	if err != nil {
		panic(fmt.Sprintf("Repository Error: %v", err))
	}
	return q
}

// Has checks if a named database pointer is registered in the contract.
func Has(name string) bool {
	regMu.RLock()
	defer regMu.RUnlock()

	q, exists := queriesMap[name]
	return exists && q != nil
}

// GetPool returns the connection pool for drivers, transactions, or tasks needing raw pool access.
// If no name is provided (or "default"), returns the default pool.
// If a specific named pool is requested but not found, returns nil.
func GetPool(name ...string) *pgxpool.Pool {
	regMu.RLock()
	defer regMu.RUnlock()

	if len(name) == 0 || name[0] == "" || name[0] == DefaultDatabasePointer {
		return pool
	}

	return poolsMap[name[0]]
}

// Reset clears all registered pools and queries (useful for test teardown).
func Reset() {
	regMu.Lock()
	defer regMu.Unlock()

	pool = nil
	queries = nil
	poolsMap = make(map[string]*pgxpool.Pool)
	queriesMap = make(map[string]*repository.Queries)
}
