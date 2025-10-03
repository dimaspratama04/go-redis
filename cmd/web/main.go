package main

import (
	"fmt"
	"golang-redis/internal/config"
	"golang-redis/internal/delivery/http/route"
	"golang-redis/internal/repository"
	"golang-redis/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	db := config.NewDatabase()

	sessionRepository := repository.NewSessionRepository()
	guestUC := usecase.NewGuestUsecase()
	authUC := usecase.NewAuthUsecase(sessionRepository)
	productUC := usecase.NewProductUseCase(db)

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
