package messaging

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func StartProductConsumer(conn *amqp.Connection) {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %v", err)
	}
	defer ch.Close()

	exchangeName := "product_events"

	err = ch.ExchangeDeclare(
		exchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
	}

	q, err := ch.QueueDeclare(
		exchangeName, // name
		true,         // durable
		false,        // autoDelete
		false,        // exclusive
		false,        // noWait
		nil,          // args
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	routingKeyTopic := "product.*"

	// queue binding
	err = ch.QueueBind(
		q.Name,          // queue name
		routingKeyTopic, // routing key
		exchangeName,    // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to bind queue: %v", err)
	}

	// consume messages
	msgs, err := ch.Consume(
		q.Name,
		routingKeyTopic,
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	log.Printf("Listening for product events on queue '%s'...", q.Name)

	// listen messages forever
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received message: %s", d.Body)
		}
	}()

	<-forever
}
