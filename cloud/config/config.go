package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	ListenPort string   `yaml:"listen_port"`
	Backends   []string `yaml:"backends"`
	RateLimit  struct {
		Capacity   int `yaml:"capacity"`
		RefillRate int `yaml:"refill_rate"`
	} `yaml:"rate_limit"`
}

func LoadConfig(path string) (*Config, error) {
	// Чтение файла конфигурации
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	// Разбор YAML
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
