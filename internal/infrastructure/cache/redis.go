package redisclient

import (
    "context"
    "log"
    "time"
    "runtime"

    "github.com/redis/go-redis/v9"
)

func NewRedisFromURL(ctx context.Context, url string) *redis.Client {
    opt, err := redis.ParseURL(url)
    if err != nil {
        log.Fatalf("invalid REDIS URL: %v", err)
    }

    // Production-friendly tuning
    // opt.PoolSize = 20                      // adjust with load; start ~20
    opt.MinIdleConns = 5
    opt.ReadTimeout = 3 * time.Second
    opt.WriteTimeout = 3 * time.Second
    opt.DialTimeout = 5 * time.Second
    opt.MaxRetries = 3
    opt.MinRetryBackoff = 100 * time.Millisecond
    opt.MaxRetryBackoff = 2 * time.Second

    // more dynamic pool sizing:
    opt.PoolSize = runtime.GOMAXPROCS(0) * 10

    rdb := redis.NewClient(opt)

    // Quick health check
    if _, err := rdb.Ping(ctx).Result(); err != nil {
        log.Fatalf("redis ping failed: %v", err)
    }

    return rdb
}

func Close(rdb *redis.Client) error {
    return rdb.Close()
}
