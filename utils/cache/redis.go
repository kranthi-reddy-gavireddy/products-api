package cache

import (
	"sync"

	"github.com/kranthi-reddy-gavireddy/products-api.git/internal/config"

	"github.com/redis/go-redis/v9"
)

var (
	cl   *redis.Client
	once sync.Once
)

func Set() error {

	var err error

	once.Do(func() {
		option, err := redis.ParseURL(config.AppConfig.RedisURL)
		if err != nil {
			return
		}

		cl = redis.NewClient(option)
	})

	return err
}

type RedisClient struct {
	*redis.Client
}

func New() *RedisClient {
	return &RedisClient{Client: cl}
}
