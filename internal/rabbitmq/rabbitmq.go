package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Komilov31/event-booker/internal/dto"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
)

const (
	retries = 3
)

type RabbitMq struct {
	pulisher *rabbitmq.Publisher
	consumer <-chan amqp.Delivery
}

func New() *RabbitMq {
	publisher, deliveries := initProducerAndConsumer()

	return &RabbitMq{
		pulisher: publisher,
		consumer: deliveries,
	}
}

func (r *RabbitMq) Publish(booking dto.QueueMessage) error {
	body, err := json.Marshal(booking)
	if err != nil {
		return fmt.Errorf("could not marshal booking to send to rabbitmq:" + err.Error())
	}

	strategy := retry.Strategy{
		Attempts: 3,
		Delay:    time.Second,
		Backoff:  2,
	}

	headers := amqp.Table{
		"x-delay": int64((time.Minute * 15).Milliseconds()), // отправляет после 15 минут(период ожидания оплаты)
	}

	options := rabbitmq.PublishingOptions{
		Headers: headers,
	}

	return r.pulisher.PublishWithRetry(body, "bookings", "application/json", strategy, options)
}

func (r *RabbitMq) Consume(ctx context.Context) (<-chan []byte, error) {
	messages := make(chan []byte)

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(messages)
				return
			default:
				next, ok := <-r.consumer
				if !ok {
					return
				}

				if err := next.Ack(false); err != nil {
					zlog.Logger.Error().Msg("could not acknowledge message consuming: " + err.Error())
				}
				messages <- next.Body
			}
		}
	}()

	return messages, nil
}
