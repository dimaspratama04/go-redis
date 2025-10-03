package controller

import (
	"golang-redis/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type GuestController struct {
	guestUC usecase.GuestUseCase
}

func NewGuestController(guestUC usecase.GuestUseCase) *GuestController {
	return &GuestController{guestUC: guestUC}
}

func (gc *GuestController) HealthCheck(c *fiber.Ctx) error {
	return gc.guestUC.HealthCheck(c)
}
