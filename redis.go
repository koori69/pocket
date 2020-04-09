// Package pocket K·J Create at 2020-04-09 21:09
package pocket

import (
	"sync"
	"time"

	redis "github.com/go-redis/redis/v7"
)

type RedisUtils struct {
	client     *redis.Client
	expiration time.Duration
	once       sync.Once
}

// RedisConfig redis config
type RedisConfig struct {
	Host       string
	Pwd        string
	ExpireTime int
}

// NewRedis get redis client
func NewRedis(config RedisConfig) (*RedisUtils, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Host,
		Password: config.Pwd, // no password set
		DB:       0,          // use default DB
	})
	_, err := redisClient.Ping().Result()
	if nil != err {
		return nil, err
	}
	return &RedisUtils{client: redisClient, expiration: time.Duration(config.ExpireTime) * time.Second}, nil
}

// Close close redis connect
func (r *RedisUtils) Close() {
	r.once.Do(func() {
		r.client.Close()
	})
}
