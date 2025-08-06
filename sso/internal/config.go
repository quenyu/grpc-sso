package internal

import (
	"errors"
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string        `yaml:"env" env-default:"local"`
	StoragePath []string      `yaml:"storage_path" env-required:"true"`
	Grpc        Grpc          `yaml:"grpc"`
	TokenTTL    time.Duration `yaml:"token_ttl"`
}

type Grpc struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func NewConfig() (*Config, error) {
	path := fetchConfigPath()
	if path == "" {
		return nil, errors.New("config path is not set")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, errors.New("config file does not exist")
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		return "./config/config.yaml"
	}
	return res
}
