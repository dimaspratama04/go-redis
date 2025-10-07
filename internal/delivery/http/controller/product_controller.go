package controller

import (
	"context"
	"encoding/json"
	"golang-redis/internal/delivery/http/request"
	"golang-redis/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ProductController struct {
	ProductUC *usecase.ProductUseCase
}

func NewProductController(productUC *usecase.ProductUseCase) *ProductController {
	return &ProductController{ProductUC: productUC}
}

func (pc *ProductController) CreateProduct(c *fiber.Ctx) error {
	body := c.Body()

	// check multiple payload
	if len(body) > 0 && body[0] == '[' {
		var reqs []request.ProductRequest
		if err := json.Unmarshal(body, &reqs); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "invalid JSON array: " + err.Error(),
			})
		}

		// use batch if multiple payload
		products, err := pc.ProductUC.CreateProductBatch(reqs)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
			})
		}

		return c.JSON(Response{
			Status:  "ok",
			Message: "product successfully created",
			Data:    products,
		})
	}

	// use single http if single payload
	var req request.ProductRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid JSON object: " + err.Error(),
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

func (pc *ProductController) GetProductByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid product id",
		})
	}

	product, err := pc.ProductUC.GetProductByID(context.Background(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "product not found",
		})
	}

	return c.JSON(Response{
		Status:  "ok",
		Message: "product fetch successfully.",
		Data:    product,
	})
}
