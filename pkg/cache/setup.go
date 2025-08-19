package cache

import (
	"time"

	"r2-challenge/cmd/envs"
)

// SetupFromEnv builds a Redis client from envs. Returns nil if REDIS_ADDR is empty.
func SetupFromEnv(environment envs.Envs) (*Client, error) {
	if environment.RedisAddr == "" {
		return nil, nil
	}

	return Setup(Config{
		Addr:       environment.RedisAddr,
		Password:   environment.RedisPassword,
		DB:         environment.RedisDB,
		DefaultTTL: 30 * time.Second,
	})
}
