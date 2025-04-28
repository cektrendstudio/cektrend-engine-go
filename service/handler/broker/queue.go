package broker

import (
	"fmt"

	"github.com/cektrendstudio/cektrend-engine-go/pkg/logger"
	"github.com/cektrendstudio/cektrend-engine-go/pkg/serror"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (ox *RabbitMQSvc) register(queue string, handler Handler) (errx serror.SError) {
	if ox.handlers == nil {
		ox.handlers = make(map[string]Handler)
	}

	if _, ok := ox.handlers[queue]; ok {
		errx = serror.Newc(fmt.Sprintf("Queue %s already exists", queue), "while registering queue")
		return errx
	}

	ox.handlers[queue] = handler
	return errx
}

func (ox *RabbitMQSvc) queues() (queues []string) {
	for queue := range ox.handlers {
		queues = append(queues, queue)
	}
	return queues
}

func (ox *RabbitMQSvc) handleMessage(queue string, msg *amqp.Delivery) {
	if handler, ok := ox.handlers[queue]; ok {
		logger.Infof("Received message from RabbitMQ queue %s", queue)

		errx := handler(msg)
		if errx != nil {
			logger.Errf("Error processing message from queue %s: %v", queue, errx)
		}
		return
	}

	logger.Warnf("Unhandled RabbitMQ queue: %s", queue)
}
