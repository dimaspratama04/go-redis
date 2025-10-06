package main

import (
	"fmt"
	"golang-redis/internal/config"
	"golang-redis/internal/delivery/http/route"
	"golang-redis/internal/repository"
	"golang-redis/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	app := fiber.New()
	db := config.NewDatabase()
	rdb := config.NewRedisClient()

	sessionRepository := repository.NewSessionRepository()
	guestUC := usecase.NewGuestUsecase()
	authUC := usecase.NewAuthUsecase(sessionRepository)
	productUC := usecase.NewProductUseCase(db, rdb)

	rc := route.RouteConfig{
		App:       app,
		AuthUC:    authUC,
		GuestUC:   guestUC,
		ProductUC: productUC}

	rc.SetupGuestRoute()
	rc.SetupAuthRoute()

	err := app.Listen(":3000")

	if err != nil {
		fmt.Sprintf("Error starting server: %s", err)
	}

}
