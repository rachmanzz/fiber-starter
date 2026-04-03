package cores

import (
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"github.com/rachmanzz/fiber-starter/config"
)

var (
	instance *config.ConfigRegistry
)

func Config() *config.ConfigRegistry {
	once.Do(func() {
		_ = godotenv.Load()
		cfg := &config.ConfigRegistry{}

		if err := env.Parse(cfg); err != nil {
			log.Fatalf("Core Config: Critical failure during parsing: %v", err)
		}

		instance = cfg
	})

	return instance
}
