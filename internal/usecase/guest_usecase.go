package usecase

import "github.com/gofiber/fiber/v2"

type GuestUseCase interface {
	HealthCheck(c *fiber.Ctx) error
}

type guestUsecase struct{}

func NewGuestUsecase() GuestUseCase {
	return &guestUsecase{}
}

func (g *guestUsecase) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "api is healthy",
	})
}
