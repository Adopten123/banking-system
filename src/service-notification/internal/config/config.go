package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string         `yaml:"env" env-default:"local"`
	RabbitMQ RabbitMQConfig `yaml:"rabbitmq"`
}

type RabbitMQConfig struct {
	URL      string `yaml:"url" env:"RABBITMQ_URL" env-required:"true"`
	Exchange string `yaml:"exchange" env-default:"transfers"`
	Queue    string `yaml:"queue" env-default:"notification_queue"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/local.yaml"
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
