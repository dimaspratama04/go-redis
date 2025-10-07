package controller

import (
	"golang-redis/internal/entity"
	"golang-redis/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	authUC usecase.AuthUsecase
}

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewAuthController(authUC usecase.AuthUsecase) *AuthController {
	return &AuthController{authUC: authUC}
}

func (ac *AuthController) Login(c *fiber.Ctx) error {
	var user entity.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid http",
		})
	}

	token, err := ac.authUC.Login(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "could not login",
		})
	}

	return c.JSON(Response{
		Status:  "ok",
		Message: "login success",
		Data: fiber.Map{
			"access_token": token,
			"user": fiber.Map{
				"id":    user.Id,
				"name":  user.Name,
				"email": user.Email,
			},
		},
	})
}

// Logout handler
func (ac *AuthController) Logout(c *fiber.Ctx) error {
	var user entity.User
	token := c.Get("Authorization")

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid http",
		})
	}

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "missing token",
		})
	}

	// ambil token dari header Bearer
	const prefix = "Bearer "
	if len(token) > len(prefix) && token[:len(prefix)] == prefix {
		token = token[len(prefix):]
	}

	if err := ac.authUC.Logout(token); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "could not logout",
		})
	}

	return c.JSON(Response{
		Status:  "ok",
		Message: "sucessfully logged out",
		Data: fiber.Map{
			"access_token": token,
			"user": fiber.Map{
				"id":    user.Id,
				"name":  user.Name,
				"email": user.Email,
			},
		},
	})

}
