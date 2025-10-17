package main

import (
	"fmt"
	"golang-redis/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	app := fiber.New()
	db := config.NewDatabase()
	rdb := config.NewRedisClient()
	rabbitConn := config.NewRabbitMQClient()

	config.Bootstrap(&config.BootstrapConfig{
		App:      app,
		Db:       db,
		Rdb:      rdb,
		RabbitMQ: rabbitConn,
	})

	err := app.Listen(":3000")

	if err != nil {
		fmt.Sprintf("Error starting server: %s", err)
	}

}
