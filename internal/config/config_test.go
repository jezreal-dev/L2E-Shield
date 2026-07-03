package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("PORT", "9090")
	os.Setenv("REDIS_URL", "redis://localhost:6379/1")
	os.Setenv("GEMINI_API_KEY", "test-api-key")
	os.Setenv("RATE_LIMIT_RPS", "12.5")
	os.Setenv("RATE_LIMIT_BURST", "20")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.Port != "9090" {
		t.Errorf("expected Port 9090, got %s", cfg.Port)
	}
	if cfg.RedisURL != "redis://localhost:6379/1" {
		t.Errorf("expected RedisURL, got %s", cfg.RedisURL)
	}
	if cfg.GeminiAPIKey != "test-api-key" {
		t.Errorf("expected GeminiAPIKey, got %s", cfg.GeminiAPIKey)
	}
	if cfg.RateLimitRPS != 12.5 {
		t.Errorf("expected RateLimitRPS 12.5, got %f", cfg.RateLimitRPS)
	}
	if cfg.RateLimitBurst != 20 {
		t.Errorf("expected RateLimitBurst 20, got %d", cfg.RateLimitBurst)
	}
}
