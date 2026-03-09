package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DB        DBConfig
	RabbitMQ  RabbitMQConfig  `yaml:"rabbitmq"`
	Exchanger ExchangerConfig `yaml:"exchanger"`
}

type DBConfig struct {
	DSN             string        `yaml:"dsn" env:"DB_URL"`
	MaxConns        int32         `yaml:"max_conns"`
	MinConns        int32         `yaml:"min_conns"`
	MaxConnIdleTime time.Duration `yaml:"max_conn_idle_time"`
	MaxConnLifetime time.Duration `yaml:"max_conn_lifetime"`
}

type RabbitMQConfig struct {
	URL string `yaml:"url" env:"RABBITMQ_URL"`
}

type ExchangerConfig struct {
	URL     string        `yaml:"url" env:"EXCHANGER_URL"`
	Timeout time.Duration `yaml:"timeout" env:"EXCHANGER_TIMEOUT" env-default:"3s"`
}

func Load(configPath string) *Config {
	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)

	if err != nil {
		log.Fatalf("cannot read config from path %s: %v", configPath, err)
	}

	return &cfg
}
