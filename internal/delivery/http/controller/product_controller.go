package controller

import (
	"golang-redis/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type Request struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type ProductController struct {
	ProductUC *usecase.ProductUseCase
}

func NewProductController(productUC *usecase.ProductUseCase) *ProductController {
	return &ProductController{ProductUC: productUC}
}

func (pc *ProductController) CreateProduct(c *fiber.Ctx) error {
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

	return c.JSON(Response{
		Status:  "ok",
		Message: "product successfully created",
		Data: fiber.Map{
			"id":           product.ID,
			"product_name": product.Name,
			"price":        product.Price,
		},
	})
}

func (pc *ProductController) GetAllProduct(c *fiber.Ctx) error {
	products, err := pc.ProductUC.GetAllProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(Response{
		Status:  "ok",
		Message: "product fetch successfully.",
		Data:    products,
	})
}
