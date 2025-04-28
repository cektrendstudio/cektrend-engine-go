package config

import (
	"fmt"

	"github.com/cektrendstudio/cektrend-engine-go/pkg/logger"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/serror"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/utils/utstring"

	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

func (cfg *Config) InitRabbitMQ() serror.SError {
	uri := fmt.Sprintf("amqp://%s:%s@%s",
		utstring.Env("BROKER_USERNAME", "guest"),
		utstring.Env("BROKER_PASSWORD", "guest"),
		utstring.Env("BROKER_ADDR", "localhost:5672"),
	)
	conn, err := amqp.Dial(uri)
	if err != nil {
		log.Error(err)
		return nil
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Error(err)
		conn.Close()
		return nil
	}

	cfg.RabbitMQConnection = conn
	cfg.RabbitMQChannel = ch

	logger.Infof("===Listening RabbitMQ on %s===", utstring.Env("BROKER_ADDR", "localhost:5672"))

	return nil
}

func (cfg *Config) Close() {
	if cfg.RabbitMQChannel != nil {
		cfg.RabbitMQChannel.Close()
	}
	if cfg.RabbitMQConnection != nil {
		cfg.RabbitMQConnection.Close()
	}
}
