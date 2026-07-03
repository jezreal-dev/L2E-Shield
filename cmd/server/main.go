package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/jezrealdev/l2e-shield/internal/cache"
	"github.com/jezrealdev/l2e-shield/internal/config"
	"github.com/jezrealdev/l2e-shield/internal/handler"
	"github.com/jezrealdev/l2e-shield/internal/middleware"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("starting L2E-Shield Go Proxy server")

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	redisCache, err := cache.New(cfg.RedisURL)
	if err != nil {
		slog.Error("failed to initialize redis cache connection client", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()
	err = redisCache.Ping(ctx)
	if err != nil {
		slog.Warn("could not establish connection ping to Redis, operating in database bypass mode", "error", err)
	} else {
		slog.Info("successfully established connection ping to Redis cache")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	proxyHandler := handler.HandleProxy(redisCache, cfg.GeminiAPIKey)
	rateLimitMiddleware := middleware.RateLimit(redisCache, cfg.RateLimitRPS, cfg.RateLimitBurst)

	mux.Handle("/v1/chat", rateLimitMiddleware(http.HandlerFunc(proxyHandler)))

	loggedRouter := middleware.Logging(mux)

	slog.Info(fmt.Sprintf("server listening on port %s", cfg.Port))
	serverErr := http.ListenAndServe(":"+cfg.Port, loggedRouter)
	if serverErr != nil {
		slog.Error("http server failed to listen", "error", serverErr)
		os.Exit(1)
	}
}
