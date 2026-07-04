// Package cache handles fast, transient data storage utilizing Redis.
package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache wraps a Redis client connection to simplify database operations.
type Cache struct {
	client *redis.Client
}

// New instantiates a new Cache object connected to the provided Redis URL.
func New(redisURL string) (*Cache, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)
	return &Cache{client: client}, nil
}

// Ping sends a heartbeat check to verify the Redis connection is alive.
func (c *Cache) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

// Get retrieves a string value from Redis by key.
func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

// Set writes a string value to Redis with an expiration time.
func (c *Cache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

// Client exposes the underlying go-redis client for advanced operations like pipelines.
func (c *Cache) Client() *redis.Client {
	return c.client
}
