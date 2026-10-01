package config

import (
	"encoding/json"
	"os"
)

// Config representa la estructura base para el MVP.
type Config struct {
	Server ServerConfig `json:"server"`
	Proxy  ProxyConfig  `json:"proxy"`
}

type ServerConfig struct {
	Port string `json:"port"`
}

type ProxyConfig struct {
	Backends []string `json:"backends"`
}

// Load lee el archivo JSON y lo deserializa.
func Load(filepath string) (*Config, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
