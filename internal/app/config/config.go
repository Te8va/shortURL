package config

import (
	"flag"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddress string
	BaseURL       string
}

func InitConfig() *Config {
	defaultServerAddress := "localhost:8080"
	defaultBaseURL := "http://localhost:8080"

	_ = godotenv.Load()

	serverAddrFlag := flag.String("a", "", "Address to run HTTP server")
	baseURLFlag := flag.String("b", "", "Base URL for short links")

	flag.Parse()

	serverAddr := getFirstNonEmpty(os.Getenv("SERVER_ADDRESS"), *serverAddrFlag, defaultServerAddress)
	baseURL := getFirstNonEmpty(os.Getenv("BASE_URL"), *baseURLFlag, defaultBaseURL)

	return &Config{
		ServerAddress: serverAddr,
		BaseURL:       baseURL,
	}
}

func getFirstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
