package config

import (
	"context"
	"strings"

	"github.com/cektrendstudio/cektrend-engine-go/pkg/serror"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/utils/utint"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/utils/utstring"

	"github.com/dgraph-io/ristretto"
	"github.com/eko/gocache/v2/cache"
	"github.com/eko/gocache/v2/store"
	redis "github.com/go-redis/redis/v8"
)

func (cfg *Config) InitCache() (errx serror.SError) {
	if utstring.Env("CACHE_ENABLED", "TRUE") != "TRUE" {
		cfg.Cache = nil
		return errx
	}

	const (
		storeMemory = "memory"
		storeRedis  = "redis"
	)

	var (
		stores    []cache.SetterCacheInterface
		storeMapx = make(map[string]store.StoreInterface)

		defaultPriorities = []string{
			storeRedis,
		}
	)

	if utstring.Env("CACHE_STORE_MEMORY_ENABLED", "TRUE") == "TRUE" {
		obj, err := ristretto.NewCache(&ristretto.Config{
			NumCounters: 1000,
			MaxCost:     100,
			BufferItems: 64,
		})
		if err != nil {
			errx = serror.NewFromErrorc(err, "Failed to create memory cache")
			return errx
		}

		storeMapx[storeMemory] = store.NewRistretto(obj, nil)
	}

	if utstring.Env("CACHE_STORE_REDIS_ENABLED", "TRUE") == "TRUE" {
		cfg.RedisClient = redis.NewClient(&redis.Options{
			Addr:     utstring.Env("CACHE_STORE_REDIS_ADDR", "127.0.0.1:6379"),
			Password: utstring.Env("CACHE_STORE_REDIS_PWD", ""),
			DB:       int(utint.StringToInt(utstring.Env("CACHE_STORE_REDIS_DB"), 0)),
		})
		storeMapx[storeRedis] = store.NewRedis(cfg.RedisClient, nil)

		GlobalShutdown.RegisterGracefullyShutdown("database/redis", func(ctx context.Context) error {
			return cfg.RedisClient.Close()
		})
	}

	priorities := utstring.CleanSpit(utstring.Env("CACHE_STORE_PRIORITIES", strings.Join(defaultPriorities, ",")), ",")
	for _, v := range priorities {
		if storeInterface, ok := storeMapx[v]; ok {
			stores = append(stores, cache.New(storeInterface))
		}
	}

	cfg.Cache = cache.NewChain(stores...)

	return errx
}
