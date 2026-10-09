package config

import (
	"net"
	"net/url"
	"strconv"
	"time"
)

type DatabaseConfig struct {
	Host            string        `env:"DB_HOST"`
	Port            int           `env:"DB_PORT" envDefault:"5432"`
	User            string        `env:"DB_USER"`
	Password        string        `env:"DB_PASSWORD"`
	Name            string        `env:"DB_NAME"`
	SSLMode         string        `env:"DB_SSLMODE" envDefault:"disable"`
	Enable          bool          `env:"DB_ENABLE" envDefault:"false"`
	MaxConns        int32         `env:"DB_MAX_CONNS" envDefault:"10"`
	MinConns        int32         `env:"DB_MIN_CONNS" envDefault:"2"`
	MaxConnLifetime time.Duration `env:"DB_MAX_CONN_LIFETIME" envDefault:"1h"`
	MaxConnIdleTime time.Duration `env:"DB_MAX_CONN_IDLE_TIME" envDefault:"30m"`
}

// DSN returns the formatted PostgreSQL connection string.
func (c DatabaseConfig) DSN() string {
	u := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(c.User, c.Password),
		Host:     net.JoinHostPort(c.Host, strconv.Itoa(c.Port)),
		Path:     c.Name,
		RawQuery: "sslmode=" + c.SSLMode,
	}
	return u.String()
}

