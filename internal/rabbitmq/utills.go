package rabbitmq

import (
	"fmt"
	"log"

	"github.com/Komilov31/event-booker/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wb-go/wbf/rabbitmq"
)

func initProducerAndConsumer() (*rabbitmq.Publisher, <-chan amqp.Delivery) {
	url := fmt.Sprintf(
		"amqp://guest:guest@%s%s/",
		config.Cfg.RabbitMq.Host,
		config.Cfg.RabbitMq.Port,
	)

	connection, err := rabbitmq.Connect(url, retries, 0)
	if err != nil {
		log.Fatal("could not connect to rabbitmq server: ", err)
	}

	pubCh, err := connection.Channel()
	if err != nil {
		log.Fatal("could not create channel for rabbitmq: ", err)
	}

	args := amqp.Table{"x-delayed-type": "direct"}
	err = pubCh.ExchangeDeclare(
		"bookings",
		"x-delayed-message",
		true,
		false,
		false,
		false,
		args,
	)
	if err != nil {
		log.Fatal("could not declare exchange for rabbitmq: ", err)
	}

	qm := rabbitmq.NewQueueManager(pubCh)
	_, err = qm.DeclareQueue("bookings")
	if err != nil {
		log.Fatal("could not create queue for rabbitmq: ", err)
	}

	err = pubCh.QueueBind("bookings", "bookings", "bookings", false, nil)
	if err != nil {
		log.Fatal("could not bind queue to exchange: ", err)
	}

	publisher := rabbitmq.NewPublisher(pubCh, "bookings")

	conCh, err := connection.Channel()
	if err != nil {
		log.Fatal("could not create channel for consumer rabbitmq: ", err)
	}

	deliveries, err := conCh.Consume("bookings", "", false, false, false, false, nil)
	if err != nil {
		log.Fatal("could not create consumer for rabbitmq: ", err)
	}

	return publisher, deliveries
}
