package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port           string
	RedisURL       string
	GeminiAPIKey   string
	RateLimitRPS   float64
	RateLimitBurst int
}

func Load() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379/0"
	}

	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	if geminiAPIKey == "" {
		return nil, fmt.Errorf("missing required environment variable: GEMINI_API_KEY")
	}

	rpsStr := os.Getenv("RATE_LIMIT_RPS")
	rps := 5.0
	if rpsStr != "" {
		val, err := strconv.ParseFloat(rpsStr, 64)
		if err == nil {
			rps = val
		}
	}

	burstStr := os.Getenv("RATE_LIMIT_BURST")
	burst := 10
	if burstStr != "" {
		val, err := strconv.Atoi(burstStr)
		if err == nil {
			burst = val
		}
	}

	return &Config{
		Port:           port,
		RedisURL:       redisURL,
		GeminiAPIKey:   geminiAPIKey,
		RateLimitRPS:   rps,
		RateLimitBurst: burst,
	}, nil
}
