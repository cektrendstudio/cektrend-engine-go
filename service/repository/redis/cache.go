package redis

import (
	"time"

	"github.com/cektrendstudio/cektrend-engine-go/pkg/logger"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/serror"
	"github.com/cektrendstudio/cektrend-engine-go/service"
	"github.com/go-redis/redis"
)

type cacheRepository struct {
	redisClient *redis.Client
}

func NewCacheRepository(redisClient *redis.Client) service.CacheRepository {
	return cacheRepository{redisClient}
}

func (c cacheRepository) Set(key string, value interface{}, duration time.Duration) error {
	if c.redisClient == nil {
		return serror.New("redis client is not ready")
	}

	return c.redisClient.Set(key, value, duration).Err()
}

func (c cacheRepository) Get(key string) (data string) {
	if c.redisClient == nil {
		logger.Info("redis client is not ready")
		return
	}

	data, err := c.redisClient.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			logger.Infof("while redis Get operation, key: %s, is not exist", key)
			return
		}

		logger.Errf("error during redis Get operation, key: %s, err: %s", key, err.Error())
		return
	}

	return
}
