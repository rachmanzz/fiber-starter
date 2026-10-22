package repocontract_test

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rachmanzz/fiber-starter/app/repository"
	"github.com/rachmanzz/fiber-starter/providers/repocontract"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContract_MultiDBRegistry(t *testing.T) {
	repocontract.Reset()
	defer repocontract.Reset()

	// 1. Uninitialized access should panic
	assert.Panics(t, func() {
		_ = repocontract.GetQueries()
	})
	assert.Panics(t, func() {
		_ = repocontract.Use()
	})
	assert.False(t, repocontract.Has("default"))
	assert.False(t, repocontract.Has("analytics"))

	// 2. Register named queries directly (mock)
	qDefault := &repository.Queries{}
	qAnalytics := &repository.Queries{}

	repocontract.RegisterNamedQueries(repocontract.DefaultDatabasePointer, qDefault)
	repocontract.RegisterNamedQueries("analytics", qAnalytics)

	// 3. Verify Has
	assert.True(t, repocontract.Has(repocontract.DefaultDatabasePointer))
	assert.True(t, repocontract.Has("analytics"))
	assert.False(t, repocontract.Has("non_existent"))

	// 4. Verify GetQueries & Use
	assert.Equal(t, qDefault, repocontract.GetQueries())
	assert.Equal(t, qDefault, repocontract.Use())
	assert.Equal(t, qDefault, repocontract.Use("default"))
	assert.Equal(t, qAnalytics, repocontract.Use("analytics"))

	// 5. Verify GetQueriesNamed
	q, err := repocontract.GetQueriesNamed("analytics")
	require.NoError(t, err)
	assert.Equal(t, qAnalytics, q)

	_, err = repocontract.GetQueriesNamed("non_existent")
	assert.Error(t, err)

	// 6. Use non-existent should panic
	assert.Panics(t, func() {
		_ = repocontract.Use("non_existent")
	})

	// 7. SetQueries overrides default
	qCustom := &repository.Queries{}
	repocontract.SetQueries(qCustom)
	assert.Equal(t, qCustom, repocontract.GetQueries())
}

func TestContract_PoolRegistry(t *testing.T) {
	repocontract.Reset()
	defer repocontract.Reset()

	dummyDefault := &pgxpool.Pool{}
	dummyEvents := &pgxpool.Pool{}

	repocontract.RegisterNamedDatabase(repocontract.DefaultDatabasePointer, dummyDefault)
	repocontract.RegisterNamedDatabase("events", dummyEvents)

	assert.Equal(t, dummyDefault, repocontract.GetPool())
	assert.Equal(t, dummyDefault, repocontract.GetPool("default"))
	assert.Equal(t, dummyEvents, repocontract.GetPool("events"))
	assert.Nil(t, repocontract.GetPool("missing"))
}
