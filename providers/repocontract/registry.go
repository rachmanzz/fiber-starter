// Package repocontract provides a centralized repository query registry for database access.
package repocontract

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
	queries    *repository.Queries
	queriesMap = make(map[string]*repository.Queries)
)

func normalizePointerName(name string) string {
	if name == "" {
		return DefaultDatabasePointer
	}
	return name
}

// DatabaseContract registers the primary PostgreSQL connection pool under the "default" pointer.
func DatabaseContract(db *pgxpool.Pool) {
	RegisterNamedDatabase(DefaultDatabasePointer, db)
}

// RegisterNamedDatabase registers a database connection pool and instantiates its SQLC Querier.
// Passing a nil pool unregisters the named database pointer.
func RegisterNamedDatabase(name string, db *pgxpool.Pool) {
	name = normalizePointerName(name)
	regMu.Lock()
	defer regMu.Unlock()

	if db == nil {
		delete(queriesMap, name)
		if name == DefaultDatabasePointer {
			queries = nil
		}
		return
	}

	q := repository.New(db)
	queriesMap[name] = q

	if name == DefaultDatabasePointer {
		queries = q
	}
}

// UnregisterNamedDatabase removes a named database pointer from the repository contract registry.
func UnregisterNamedDatabase(name string) {
	RegisterNamedDatabase(name, nil)
}

// RegisterNamedQueries registers a pre-instantiated *repository.Queries under a named pointer (useful for testing and mocks).
func RegisterNamedQueries(name string, q *repository.Queries) {
	name = normalizePointerName(name)
	regMu.Lock()
	defer regMu.Unlock()

	if q == nil {
		delete(queriesMap, name)
		if name == DefaultDatabasePointer {
			queries = nil
		}
		return
	}

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
	target := normalizePointerName(name)
	regMu.RLock()
	defer regMu.RUnlock()

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
	if len(name) > 0 {
		target = normalizePointerName(name[0])
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
	target := normalizePointerName(name)
	regMu.RLock()
	defer regMu.RUnlock()

	q, exists := queriesMap[target]
	return exists && q != nil
}

// Reset clears all registered queries (useful for test teardown).
func Reset() {
	regMu.Lock()
	defer regMu.Unlock()

	queries = nil
	queriesMap = make(map[string]*repository.Queries)
}
