package cache

import (
    "context"
    "time"

    "github.com/redis/go-redis/v9"
)

type Client struct {
    Redis *redis.Client
}

type Config struct {
    Addr       string
    Password   string
    DB         int
    DefaultTTL time.Duration
}

func Setup(cfg Config) (*Client, error) {
    redisClient := redis.NewClient(&redis.Options{
        Addr:     cfg.Addr,
        Password: cfg.Password,
        DB:       cfg.DB,
    })

    if err := redisClient.Ping(context.Background()).Err(); err != nil {
        return nil, err
    }

    return &Client{Redis: redisClient}, nil
}

func (c *Client) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
    return c.Redis.Set(ctx, key, value, ttl).Err()
}

func (c *Client) Get(ctx context.Context, key string) ([]byte, error) {
    data, err := c.Redis.Get(ctx, key).Bytes()
    if err == redis.Nil {
        return nil, nil
    }

    return data, err
}

func (c *Client) Del(ctx context.Context, keys ...string) error {
    return c.Redis.Del(ctx, keys...).Err()
}


