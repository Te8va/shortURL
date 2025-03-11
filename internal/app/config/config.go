package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddress    string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL          string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath  string `env:"FILE_STORAGE_PATH" envDefault:"storage.json"`
	PostgresUser     string `env:"POSTGRES_USER"         envDefault:"shortURL"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"     envDefault:"shortURL"`
	PostgresDB       string `env:"POSTGRES_DB"     envDefault:"shortURL"`
	PostgresPort     int    `env:"POSTGRES_PORT"         envDefault:"5432"`
	PostgresConn     string `env:"POSTGRES_CONN"  envDefault:"postgres://shortURL:shortURL@shortURL-db:5432/shortURL?sslmode=disable"`
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
	if *databaseDSNFlag != "" {
		cfg.PostgresConn = *databaseDSNFlag
	}

	return &cfg
}
