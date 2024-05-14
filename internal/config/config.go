package config

import "github.com/caarlos0/env/v11"

type Config struct {
	Server
	Postgres
}

func NewConfig() (*Config, error) {
	cfg := new(Config)

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
