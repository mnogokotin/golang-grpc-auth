package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type (
	Config struct {
		Env      string        `yaml:"env"`
		Grpc     Grpc          `yaml:"grpc"`
		Pg       Pg            `yaml:"pg"`
		TokenTtl time.Duration `yaml:"token_ttl"`
	}

	Grpc struct {
		Port    int           `yaml:"port"`
		Timeout time.Duration `yaml:"timeout"`
	}

	Pg struct {
		Url string `yaml:"url" env:"PG_URL"`
	}
)

func New() *Config {
	configPath := "./config/main.yml" // os.Getenv("CONFIG_PATH")
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("can't read config: " + err.Error())
	}

	return &cfg
}
