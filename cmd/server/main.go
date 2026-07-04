package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	corsMiddleware := middleware.CORS(cfg.AllowedOrigin)

	mux.Handle("/v1/chat", rateLimitMiddleware(http.HandlerFunc(proxyHandler)))

	// Stack middlewares: Logging -> CORS -> Multiplexer
	finalRouter := middleware.Logging(corsMiddleware(mux))

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      finalRouter,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Server shutdown channel
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info(fmt.Sprintf("server listening on port %s", cfg.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("http server failed to listen", "error", err)
			os.Exit(1)
		}
	}()

	<-stop
	slog.Info("shutting down gracefully, pressing Ctrl+C again will force exit")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("server exited cleanly")
}
