package main

import (
	"golang-redis/internal/config"
	"golang-redis/internal/delivery/messaging"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	conn := config.NewRabbitMQClient()
	defer conn.Close()

	log.Println("Starting Product Consumer Worker...")
	messaging.StartProductConsumer(conn)
}
