package cores_test

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rachmanzz/fiber-starter/cores"
	"github.com/stretchr/testify/assert"
)

func TestNamedDatabasePools_Registry(t *testing.T) {
	// Clean slate
	cores.CloseDB()
	defer cores.CloseDB()

	// 1. Should report not found when empty
	assert.False(t, cores.HasDB("custom_pool"))

	// 2. Dummy pool reference
	dummyPool := &pgxpool.Pool{}

	cores.RegisterDBPool("event", dummyPool)
	cores.RegisterDBPool(cores.DefaultDatabaseName, dummyPool)

	// 3. Verify HasDB
	assert.True(t, cores.HasDB("event"))
	assert.True(t, cores.HasDB(cores.DefaultDatabaseName))
	assert.True(t, cores.HasDB("default"))
	assert.False(t, cores.HasDB("non_existent"))

	// 4. CloseNamedDB
	err := cores.CloseNamedDB("event")
	assert.NoError(t, err)
	assert.False(t, cores.HasDB("event"))
	assert.True(t, cores.HasDB("default"))

	// 5. CloseNamedDB on non-existent pool
	err = cores.CloseNamedDB("event")
	assert.Error(t, err)

	// 6. CloseDB closes all
	cores.CloseDB()
	assert.False(t, cores.HasDB("default"))
}

func TestNamedDatabaseContract_CallbackAndDeadlockFree(t *testing.T) {
	cores.CloseDB()
	defer cores.CloseDB()

	received := make(map[string]*pgxpool.Pool)
	// Verify that the callback can safely invoke cores.HasDB without deadlocking
	cores.SetNamedDatabaseContract(func(name string, pool *pgxpool.Pool) {
		received[name] = pool
		_ = cores.HasDB(name)
	})

	dummyPool := &pgxpool.Pool{}
	cores.RegisterDBPool("analytics", dummyPool)

	assert.Equal(t, dummyPool, received["analytics"])
}

func TestNamedDatabaseContract_ExistingPoolsBackfill(t *testing.T) {
	cores.CloseDB()
	defer cores.CloseDB()

	dummyPool := &pgxpool.Pool{}
	cores.RegisterDBPool("existing_pool", dummyPool)

	backfilled := make(map[string]*pgxpool.Pool)
	// Even though SetNamedDatabaseContract is called after RegisterDBPool, it should immediately backfill
	cores.SetNamedDatabaseContract(func(name string, pool *pgxpool.Pool) {
		backfilled[name] = pool
	})

	assert.Equal(t, dummyPool, backfilled["existing_pool"])
}

func TestDatabaseContract_DefaultPoolBackfill(t *testing.T) {
	cores.CloseDB()
	defer cores.CloseDB()

	dummyPool := &pgxpool.Pool{}
	cores.RegisterDBPool(cores.DefaultDatabaseName, dummyPool)

	var defaultPool *pgxpool.Pool
	cores.SetDatabaseContract(func(pool *pgxpool.Pool) {
		defaultPool = pool
	})

	assert.Equal(t, dummyPool, defaultPool)
}
