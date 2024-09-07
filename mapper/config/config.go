package config

import "github.com/caarlos0/env/v11"

type Config struct {
	APIToken       string   `env:"API_TOKEN" envDefault:"llhls-token"`
	APIHosts       []string `env:"API_HOSTS" envSeparator:","`
	RedisAddr      string   `env:"REDIS_ADDR" envDefault:"localhost:6379"`
	StreamTTL      int      `env:"STREAM_TTL" envDefault:"10"`
	UpdateInterval int      `env:"UPDATE_INTERVAL" envDefault:"5"`
}

func New() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return &cfg, err
	}

	return &cfg, nil
}
