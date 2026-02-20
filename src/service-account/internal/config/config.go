package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DB DBConfig
}

type DBConfig struct {
	DSN             string        `yaml:"dsn" env:"DB_URL"`
	MaxConns        int32         `yaml:"max_conns"`
	MinConns        int32         `yaml:"min_conns"`
	MaxConnIdleTime time.Duration `yaml:"max_conn_idle_time"`
	MaxConnLifetime time.Duration `yaml:"max_conn_lifetime"`
}

func Load(configPath string) *Config {
	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)

	if err != nil {
		log.Fatalf("cannot read config from path %s: %v", configPath, err)
	}

	return &cfg
}
