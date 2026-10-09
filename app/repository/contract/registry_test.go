package contract_test

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rachmanzz/fiber-starter/app/repository"
	"github.com/rachmanzz/fiber-starter/app/repository/contract"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContract_MultiDBRegistry(t *testing.T) {
	contract.Reset()
	defer contract.Reset()

	// 1. Uninitialized access should panic
	assert.Panics(t, func() {
		_ = contract.GetQueries()
	})
	assert.Panics(t, func() {
		_ = contract.Use()
	})
	assert.False(t, contract.Has("default"))
	assert.False(t, contract.Has("analytics"))

	// 2. Register named queries directly (mock)
	qDefault := &repository.Queries{}
	qAnalytics := &repository.Queries{}

	contract.RegisterNamedQueries(contract.DefaultDatabasePointer, qDefault)
	contract.RegisterNamedQueries("analytics", qAnalytics)

	// 3. Verify Has
	assert.True(t, contract.Has(contract.DefaultDatabasePointer))
	assert.True(t, contract.Has("analytics"))
	assert.False(t, contract.Has("non_existent"))

	// 4. Verify GetQueries & Use
	assert.Equal(t, qDefault, contract.GetQueries())
	assert.Equal(t, qDefault, contract.Use())
	assert.Equal(t, qDefault, contract.Use("default"))
	assert.Equal(t, qAnalytics, contract.Use("analytics"))

	// 5. Verify GetQueriesNamed
	q, err := contract.GetQueriesNamed("analytics")
	require.NoError(t, err)
	assert.Equal(t, qAnalytics, q)

	_, err = contract.GetQueriesNamed("non_existent")
	assert.Error(t, err)

	// 6. Use non-existent should panic
	assert.Panics(t, func() {
		_ = contract.Use("non_existent")
	})

	// 7. SetQueries overrides default
	qCustom := &repository.Queries{}
	contract.SetQueries(qCustom)
	assert.Equal(t, qCustom, contract.GetQueries())
}

func TestContract_PoolRegistry(t *testing.T) {
	contract.Reset()
	defer contract.Reset()

	dummyDefault := &pgxpool.Pool{}
	dummyEvents := &pgxpool.Pool{}

	contract.RegisterNamedDatabase(contract.DefaultDatabasePointer, dummyDefault)
	contract.RegisterNamedDatabase("events", dummyEvents)

	assert.Equal(t, dummyDefault, contract.GetPool())
	assert.Equal(t, dummyDefault, contract.GetPool("default"))
	assert.Equal(t, dummyEvents, contract.GetPool("events"))
	assert.Nil(t, contract.GetPool("missing"))
}
