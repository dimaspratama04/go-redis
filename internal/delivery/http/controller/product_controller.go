package controller

import (
	"context"
	"encoding/json"
	"golang-redis/internal/delivery/http/request"
	"golang-redis/internal/delivery/http/response"
	"golang-redis/internal/usecase"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

type ProductController struct {
	ProductUC *usecase.ProductUseCase
}

func NewProductController(productUC *usecase.ProductUseCase) *ProductController {
	return &ProductController{ProductUC: productUC}
}

func (pc *ProductController) GetProducts(c *fiber.Ctx) error {
	queryParams := new(request.ProductRequestQueryParams)

	// error parser queryparams
	if err := c.QueryParser(queryParams); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid query format.",
		})
	}

	// error query invalid / not exist
	if err := queryParams.ValidateQuery(c); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid query parameter.",
		})
	}

	// get by category
	category := queryParams.Category
	if category != "" {

		products, err := pc.ProductUC.GetProductByCategory(context.Background(), category)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "failed to fetch products by category: ",
			})
		}

		if len(products) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "no product found for category: " + category,
			})
		}

		return c.JSON(Response{
			Status:  "ok",
			Message: "products filtered by category successfully.",
			Data:    products,
		})
	}

	// get all products
	ctx := context.Background()
	allProducts, err := pc.ProductUC.GetAllProducts(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to fetch all products",
		})
	}

	return c.JSON(Response{
		Status:  "ok",
		Message: "all products fetched successfully.",
		Data:    allProducts,
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

	response := response.ProductResponse{
		ID:       uint64(product.ID),
		Name:     product.Name,
		Price:    product.Price,
		Category: product.Category,
	}

	return c.JSON(Response{
		Status:  "ok",
		Message: "product successfully created",
		Data:    response,
	})
}

func (pc *ProductController) CreateProduct(c *fiber.Ctx) error {
	body := c.Body()

	// multiple entries
	if len(body) > 0 && body[0] == '[' {
		var reqs []request.ProductRequest
		if err := json.Unmarshal(body, &reqs); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "invalid JSON array: " + err.Error(),
			})
		}

		for _, req := range reqs {
			if err := validate.Struct(req); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status":  "error",
					"message": "Validation failed: One or more product entries in the batch are missing required keys.",
				})
			}
		}

		products, err := pc.ProductUC.CreateProductBatch(context.Background(), reqs)
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

	// single entries
	var req request.ProductRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid JSON object: " + err.Error(),
		})
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Validation failed: are missing required keys.",
		})
	}

	product, err := pc.ProductUC.CreateProduct(context.Background(), req.Name, req.Price, req.Category)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	response := response.ProductResponse{
		ID:       uint64(product.ID),
		Name:     product.Name,
		Price:    product.Price,
		Category: product.Category,
	}

	return c.JSON(Response{
		Status:  "ok",
		Message: "product successfully created",
		Data:    response,
	})
}
