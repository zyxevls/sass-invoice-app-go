package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func NewRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}

// CacheGet retrieves a cached value as JSON unmarshaled to the provided interface
func CacheGet(rdb *redis.Client, key string, dest interface{}) error {
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

// CacheSet stores a value as JSON with expiration time
func CacheSet(rdb *redis.Client, key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, key, jsonData, expiration).Err()
}

// CacheInvalidate deletes keys matching a pattern
func CacheInvalidate(rdb *redis.Client, pattern string) error {
	iter := rdb.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		rdb.Del(ctx, key)
	}
	return iter.Err()
}

func RateLimiter(rdb *redis.Client) fiber.Handler {
	return func(c fiber.Ctx) error {
		ip := c.IP()
		key := "rate:" + ip

		count, _ := rdb.Incr(ctx, key).Result()

		if count == 1 {
			rdb.Expire(ctx, key, time.Minute)
		}

		if count > 100 {
			return c.Status(429).JSON(fiber.Map{"error": "Too many requests"})
		}

		return c.Next()
	}
}
