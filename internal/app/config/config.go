package config

import (
	"flag"
	"log"
	"os"

	"github.com/caarlos0/env/v6"
)

// Config holds application configuration parameters
type Config struct {
	ServerAddress    string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL          string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath  string `env:"FILE_STORAGE_PATH"`
	PostgresUser     string `env:"POSTGRES_USER"         envDefault:"shortURL"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"     envDefault:"shortURL"`
	PostgresDB       string `env:"POSTGRES_DB"     envDefault:"shortURL"`
	PostgresPort     int    `env:"POSTGRES_PORT"         envDefault:"5432"`
	DatabaseDSN      string `env:"DATABASE_DSN"`
	JWTKey           string `env:"JWT_KEY"               envDefault:"supermegasecret"`
}

// NewConfig creates and returns a Config instance by parsing environment variables and command-line flags.
func NewConfig() *Config {

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Error parsing environment variables:%v", err)
	}

	serverAddrFlag := flag.String("a", "", "Address to run HTTP server")
	baseURLFlag := flag.String("b", "", "Base URL for short links")
	fileStorageFlag := flag.String("f", "", "Path to storage file")
	databaseDSNFlag := flag.String("d", "", "PostgreSQL connection string")

	flag.Parse()

	if *serverAddrFlag != "" {
		cfg.ServerAddress = *serverAddrFlag
	}
	if *baseURLFlag != "" {
		cfg.BaseURL = *baseURLFlag
	}
	if *fileStorageFlag != "" {
		cfg.FileStoragePath = *fileStorageFlag
	}

	if dsnEnv, exists := os.LookupEnv("DATABASE_DSN"); exists && dsnEnv != "" {
		cfg.DatabaseDSN = dsnEnv
	} else if *databaseDSNFlag != "" {
		cfg.DatabaseDSN = *databaseDSNFlag
	}

	return &cfg
}
