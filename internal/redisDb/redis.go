package redisDb

import (
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vickon16/go-gin-rest-api/internal/env"
	"github.com/vickon16/go-gin-rest-api/internal/utils"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient() *RedisClient {
	address := env.GetEnvString("REDIS_ADDRESS", "localhost:6379")

	redisDb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       0, // use default DB
	})

	ctx, cancel := utils.CreateContext()
	defer cancel()

	// Test connection
	_, err := redisDb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	log.Println("Redis connect successfully")
	return &RedisClient{Client: redisDb}
}

func (r *RedisClient) Set(key string, value any, duration time.Duration) error {
	ctx, cancel := utils.CreateContext()
	defer cancel()

	return r.Client.Set(ctx, key, value, duration).Err()
}

func (r *RedisClient) Get(key string) (string, error) {
	ctx, cancel := utils.CreateContext()
	defer cancel()

	return r.Client.Get(ctx, key).Result()
}

func (r *RedisClient) Delete(key string) error {
	ctx, cancel := utils.CreateContext()
	defer cancel()

	return r.Client.Del(ctx, key).Err()
}
