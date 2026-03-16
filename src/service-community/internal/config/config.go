package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Env      string         `yaml:"env"`
	HTTP     HTTPConfig     `yaml:"http"`
	DB       DBConfig       `yaml:"db"`
	Redis    RedisConfig    `yaml:"redis"`
	RabbitMQ RabbitMQConfig `yaml:"rabbitmq"`
}

type HTTPConfig struct {
	Port int `yaml:"port"`
}

type DBConfig struct {
	URL string `yaml:"url"`
}

type RedisConfig struct {
	Addr string `yaml:"addr"`
}

type RabbitMQConfig struct {
	URL string `yaml:"url"`
}

// Load read Config
func Load() (*Config, error) {

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/local.yaml"
	}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return &cfg, nil
}
