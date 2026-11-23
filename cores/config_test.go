package cores_test

import (
	"testing"

	"github.com/rachmanzz/fiber-starter/cores"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadAndReloadConfig(t *testing.T) {
	t.Setenv("APP_NAME", "CustomReloadedApp")
	t.Setenv("DB_HOST", "db.internal")

	cfg, err := cores.ReloadConfig()
	require.NoError(t, err)
	assert.Equal(t, "CustomReloadedApp", cfg.App.Name)
	assert.Equal(t, "db.internal", cfg.Database.Host)
	assert.Equal(t, "CustomReloadedApp", cores.Config().App.Name)
}
