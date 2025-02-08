package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error to load .env file or not found:%v", err)
	}

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Error parsing environment variables:%v", err)
	}

	serverAddrFlag := flag.String("a", "", "Address to run HTTP server")
	baseURLFlag := flag.String("b", "", "Base URL for short links")

	flag.Parse()

	if *serverAddrFlag != "" {
		cfg.ServerAddress = *serverAddrFlag
	}
	if *baseURLFlag != "" {
		cfg.BaseURL = *baseURLFlag
	}

	return &cfg
}
