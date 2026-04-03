package config

type LoggerConfig struct {
	Path       string `env:"LOGGER_PATH" envDefault:"logs/app.log"`
	MaxSize    int    `env:"LOG_MAX_SIZE" envDefault:"100"`
	MaxBackups int    `env:"LOG_MAX_BACKUPS" envDefault:"3"`
	MaxAge     int    `env:"LOG_MAX_AGE" envDefault:"28"`
	Compress   bool   `env:"LOG_COMPRESS" envDefault:"true"`
	Level      string `env:"LOG_LEVEL" envDefault:"debug"`
}
