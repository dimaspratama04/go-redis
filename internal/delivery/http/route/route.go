package route

import (
	"golang-redis/internal/delivery/http/controller"
	"golang-redis/internal/usecase"

	"strings"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App       *fiber.App
	AuthUC    usecase.AuthUsecase
	GuestUC   usecase.GuestUseCase
	ProductUC *usecase.ProductUseCase
}

func (rc *RouteConfig) AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "missing authorization header",
		})
	}

	// format: Bearer <token>
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid token format",
		})
	}

	token := parts[1]
	session, _ := rc.AuthUC.ValidateToken(token)
	if session == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid or expired token",
		})
	}

	// simpan user id di context
	c.Locals("userID", session.UserID)
	return c.Next()
}

// Guest route (tidak butuh login)
func (rc *RouteConfig) SetupGuestRoute() {

	authController := controller.NewAuthController(rc.AuthUC)
	guestController := controller.NewGuestController(rc.GuestUC)

	rc.App.Get("/api/healthz", guestController.HealthCheck)
	rc.App.Post("/api/users/login", authController.Login)
}

// Auth route (butuh login)
func (rc *RouteConfig) SetupAuthRoute() {
	authController := controller.NewAuthController(rc.AuthUC)
	productController := controller.NewProductController(rc.ProductUC)

	rc.App.Use(rc.AuthMiddleware)

	rc.App.Post("/api/users/logout", authController.Logout)
	rc.App.Get("/api/products", productController.List)
	rc.App.Post("/api/products", productController.Create)
}
