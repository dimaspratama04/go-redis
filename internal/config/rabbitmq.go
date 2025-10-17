package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rabbitmq/amqp091-go"
)

func NewRabbitMQClient() *amqp091.Connection {
	host := os.Getenv("RABBITMQ_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("RABBITMQ_PORT")
	if port == "" {
		port = "5672"
	}

	user := os.Getenv("RABBITMQ_USER")
	if user == "" {
		user = "guest"
	}

	pass := os.Getenv("RABBITMQ_PASS")
	if pass == "" {
		pass = "guest"
	}

	vhost := os.Getenv("RABBITMQ_VHOST")
	if vhost == "" {
		vhost = "/"
	}

	if vhost != "/" && strings.HasPrefix(vhost, "/") {
		vhost = vhost[1:]
	}

	var rabbitURL string
	if strings.Contains(host, ":") {
		rabbitURL = fmt.Sprintf("amqp://%s:%s@%s/%s", user, pass, host, vhost)
	} else {
		rabbitURL = fmt.Sprintf("amqp://%s:%s@%s:%s/%s", user, pass, host, port, vhost)
	}

	conn, err := amqp091.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to RabbitMQ: %v", err)
	}

	log.Printf("✅ Connected to RabbitMQ at %s", rabbitURL)
	return conn
}
