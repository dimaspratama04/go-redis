package controller

import (
	"golang-redis/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type ProductController struct {
	ProductUC *usecase.ProductUseCase
}

func NewProductController(productUC *usecase.ProductUseCase) *ProductController {
	return &ProductController{ProductUC: productUC}
}

func (pc *ProductController) Create(c *fiber.Ctx) error {
	type Request struct {
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}

	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	product, err := pc.ProductUC.CreateProduct(req.Name, req.Price)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "product created",
		"data":    product,
	})
}

func (pc *ProductController) List(c *fiber.Ctx) error {
	products, err := pc.ProductUC.GetAllProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "success",
		"data":    products,
	})
}
