package config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"

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
	EnableHTTPS      bool
}

// ConfigFile describes JSON configuration file format
type ConfigFile struct {
	ServerAddress   string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDSN     string `json:"database_dsn"`
	EnableHTTPS     bool   `json:"enable_https"`
}

func loadFromFile(path string) (*ConfigFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfgFile ConfigFile
	if err := json.NewDecoder(f).Decode(&cfgFile); err != nil {
		return nil, err
	}
	return &cfgFile, nil
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
	httpsFlag := flag.Bool("s", false, "Enable HTTPS")
	configPathFlag := flag.String("c", "", "Path to config file (JSON)")
	configPathFlagLong := flag.String("config", "", "Path to config file (JSON)")

	flag.Parse()

	configPath := ""
	if *configPathFlag != "" {
		configPath = *configPathFlag
	} else if *configPathFlagLong != "" {
		configPath = *configPathFlagLong
	} else if envPath := os.Getenv("CONFIG"); envPath != "" {
		configPath = envPath
	}

	if configPath != "" {
		if cfgFile, err := loadFromFile(configPath); err == nil {
			cfg.ServerAddress = cfgFile.ServerAddress
			cfg.BaseURL = cfgFile.BaseURL
			cfg.FileStoragePath = cfgFile.FileStoragePath
			cfg.DatabaseDSN = cfgFile.DatabaseDSN
			cfg.EnableHTTPS = cfgFile.EnableHTTPS
		}
	}

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

	cfg.EnableHTTPS = *httpsFlag
	if secureEnv, exists := os.LookupEnv("ENABLE_HTTPS"); exists || !cfg.EnableHTTPS {
		if val, err := strconv.ParseBool(secureEnv); err == nil {
			cfg.EnableHTTPS = val
		}
	}

	return &cfg
}
