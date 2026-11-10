package config

import (
	"os"

	"go.yaml.in/yaml/v3"
)

type Config struct {
	DATABASE DatabaseConfig `yaml:"database"`
	SERVER   ServerConfig   `yaml:"server"`
}

type DatabaseConfig struct {
	HOSTNAME string `yaml:"hostname"`
	USERNAME string `yaml:"username"`
	PASSWORD string `yaml:"password"`
	DATABASE string `yaml:"database"`
	PORT     int    `yaml:"port"`
}

type ServerConfig struct {
	HOSTNAME string `yaml:"hostname"`
	PORT     int    `yaml:"port"`
}

func LoadConfig(path string) (Config, error) {
	var cfg Config
	if path == "" {
		path = "config.yaml"
	}
	file, err := os.ReadFile(path)

	if err != nil {
		return cfg, err
	}

	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
