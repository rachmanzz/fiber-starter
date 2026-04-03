package config

type AppConfig struct {
	Name    string `env:"APP_NAME" envDefault:"FiberApp"`
	Port    string `env:"APP_PORT" envDefault:":3000"`
	Env     string `env:"APP_ENV" envDefault:"development"`
	Debug   bool   `env:"APP_DEBUG" envDefault:"false"`
	Version string `env:"APP_VERSION" envDefault:"1.0.0"`
}
