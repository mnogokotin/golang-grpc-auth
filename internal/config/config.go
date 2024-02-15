package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	Env            string `yaml:"env"`
	StoragePath    string `yaml:"storage_path" env-required:"true"`
	Grpc           Grpc   `yaml:"grpc"`
	MigrationsPath string
	TokenTtl       time.Duration `yaml:"token_ttl"`
}

type Grpc struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func New() *Config {
	configPath := "./config/main.yml" // os.Getenv("CONFIG_PATH")
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}
