package broker

import (
	"context"
	"fmt"

	"github.com/cektrendstudio/cektrend-engine-go/pkg/logger"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/serror"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Handler func(msg *amqp.Delivery) (errx serror.SError)

type RabbitMQSvc struct {
	isAutoStarted bool
	channel       *amqp.Channel
	handlers      map[string]Handler
}

func NewRabbitMQHandler(
	isAutoStarted bool,
	channel *amqp.Channel,
) *RabbitMQSvc {
	obj := RabbitMQSvc{
		isAutoStarted: isAutoStarted,
		channel:       channel,
	}

	// obj.register(models.BrokerQueueTransactionNotification, obj.handlerTransactionNotification)

	if obj.channel == nil {
		logger.Errf("RabbitMQ channel is nil")
		return nil
	}

	if isAutoStarted {
		go obj.ConsumeMessages(context.Background())
	}

	return &obj
}

func (ox *RabbitMQSvc) ConsumeMessages(ctx context.Context) error {
	queues := ox.queues()
	for _, queue := range queues {
		msgs, err := ox.channel.Consume(
			queue,
			"",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			logger.Errf("Error while subscribing to queue %s", queue)
			return serror.NewFromErrorc(err, fmt.Sprintf("while subscribing to queue %s", queue))
		}

		go func(queue string, msgs <-chan amqp.Delivery) {
			for {
				select {
				case <-ctx.Done():
					logger.Infof("Stopping consumption for queue: %s", queue)
					return
				case msg := <-msgs:
					if msg.Body == nil {
						logger.Warnf("Empty message received from queue: %s", queue)
						continue
					}
					logger.Infof("Message received from queue %s: %s", queue, string(msg.Body))
					ox.handleMessage(queue, &msg)
				}
			}
		}(queue, msgs)
	}

	return nil
}
