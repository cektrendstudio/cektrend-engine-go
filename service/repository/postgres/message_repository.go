package postgres

import (
	"encoding/json"

	"github.com/cektrendstudio/cektrend-engine-go/pkg/logger"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/serror"
	"github.com/cektrendstudio/cektrend-engine-go/service"

	amqp "github.com/rabbitmq/amqp091-go"
)

type messageRepo struct {
	channel *amqp.Channel
}

func NewMessageRepository(
	channel *amqp.Channel,
) service.MessageRepository {
	return &messageRepo{
		channel: channel,
	}
}

func (r *messageRepo) Publish(queue string, data interface{}) (errx serror.SError) {
	byt, err := json.Marshal(data)
	if err != nil {
		errx = serror.NewFromErrorc(err, "while marshaling data")
		return errx
	}

	err = r.channel.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        byt,
		},
	)
	if err != nil {
		errx = serror.NewFromErrorc(err, "while publishing to RabbitMQ")
		return errx
	}

	logger.Infof("Message successfully published to queue: %s", queue)
	return errx
}
