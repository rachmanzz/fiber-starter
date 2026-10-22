package config_test

import (
	"testing"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/rachmanzz/fiber-starter/config"
)

func TestDatabaseConfig_DurationParsing(t *testing.T) {
	t.Setenv("DB_MAX_CONN_LIFETIME", "1h30m")
	t.Setenv("DB_MAX_CONN_IDLE_TIME", "15m")

	cfg := &config.ConfigRegistry{}
	err := env.Parse(cfg)
	if err != nil {
		t.Fatalf("Failed to parse config: %v", err)
	}

	expectedLifetime := 1*time.Hour + 30*time.Minute
	if cfg.Database.MaxConnLifetime != expectedLifetime {
		t.Errorf("Expected MaxConnLifetime %v, got %v", expectedLifetime, cfg.Database.MaxConnLifetime)
	}

	expectedIdleTime := 15 * time.Minute
	if cfg.Database.MaxConnIdleTime != expectedIdleTime {
		t.Errorf("Expected MaxConnIdleTime %v, got %v", expectedIdleTime, cfg.Database.MaxConnIdleTime)
	}
}
