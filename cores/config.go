package cores

import (
	"log"
	"sync"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"github.com/rachmanzz/fiber-starter/config"
)

var (
	configMu sync.RWMutex
	instance *config.ConfigRegistry
)

// LoadConfig parses environment variables into a new ConfigRegistry without terminating the process on failure.
func LoadConfig() (*config.ConfigRegistry, error) {
	_ = godotenv.Load()
	cfg := &config.ConfigRegistry{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// ReloadConfig re-parses environment variables and updates the active ConfigRegistry singleton.
func ReloadConfig() (*config.ConfigRegistry, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}
	configMu.Lock()
	instance = cfg
	configMu.Unlock()
	return cfg, nil
}

func Config() *config.ConfigRegistry {
	configMu.RLock()
	if instance != nil {
		cfg := instance
		configMu.RUnlock()
		return cfg
	}
	configMu.RUnlock()

	configMu.Lock()
	defer configMu.Unlock()
	if instance == nil {
		cfg, err := LoadConfig()
		if err != nil {
			log.Fatalf("Core Config: Critical failure during parsing: %v", err)
		}
		instance = cfg
	}
	return instance
}
