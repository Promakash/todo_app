package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
	"todo/pkg/infra/cache"
)

const CacheAlwaysAlive = redis.KeepTTL

type Redis struct {
	client *redis.Client
}

func NewRedisClient(cfg Config) (*redis.Client, error) {
	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: cfg.Password,
		DB:       0,
		Protocol: 3,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}

func NewRedisService(client *redis.Client) cache.Cache {
	return &Redis{client: client}
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}, TTL time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, bytes, TTL).Err()
}

func (r *Redis) Get(ctx context.Context, key string, value interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), value)
}

func (r *Redis) Delete(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

func ShutdownClient(client *redis.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = client.Shutdown(ctx)
}
