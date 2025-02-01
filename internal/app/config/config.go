package config

import "flag"

type Config struct {
	ServerAddress string
	BaseURL       string
}

func InitConfig() *Config {
	cfg := &Config{}

	flag.StringVar(&cfg.ServerAddress, "a", "localhost:8080", "Address to run HTTP server")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "Base URL for short links")

	flag.Parse()
	return cfg
}
