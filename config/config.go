package config

import (
	"os"

	"github.com/cektrendstudio/cektrend-engine-go/models"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/serror"

	"github.com/eko/gocache/v2/cache"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	DB                 *sqlx.DB
	Server             *gin.Engine
	RedisClient        *redis.Client
	Cache              *cache.ChainCache
	AWSConfig          *models.AWSConfig
	RabbitMQConnection *amqp.Connection
	RabbitMQChannel    *amqp.Channel
}

func Init() (cfg Config) {
	Catch(cfg.InitTimezone())
	Catch(cfg.InitBucket())
	Catch(cfg.InitPostgres())
	Catch(cfg.InitRabbitMQ())
	Catch(cfg.InitCache())
	Catch(cfg.InitService())

	return
}

func (cfg *Config) Start() (errx serror.SError) {
	cfg.Server.Run(os.Getenv("APP_PORT"))

	return
}

func Catch(errx serror.SError) {
	if errx != nil {
		errx.Panic()
	}
}
