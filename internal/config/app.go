package config

import (
	"golang-redis/internal/delivery/http/route"
	"golang-redis/internal/repository"
	"golang-redis/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	Db       *gorm.DB
	App      *fiber.App
	Rdb      *redis.Client
	RabbitMQ *amqp091.Connection
}

func Bootstrap(config *BootstrapConfig) {
	// setup repository
	sessionRepository := repository.NewSessionRepository()
	productRepository := repository.NewProductRepository(config.Db, config.Rdb)

	// setup usecase
	guestUC := usecase.NewGuestUsecase()
	authUC := usecase.NewAuthUsecase(sessionRepository)
	productUC := usecase.NewProductUseCase(productRepository)

	// router config
	rc := route.RouteConfig{
		App:       config.App,
		AuthUC:    authUC,
		GuestUC:   guestUC,
		ProductUC: productUC}

	rc.SetupGuestRoute()
	rc.SetupAuthRoute()
}
