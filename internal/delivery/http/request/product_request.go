package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ProductRequestQueryParams struct {
	Category string `query:"category"`
}

type ProductRequest struct {
	Name     string  `json:"name" validate:"required"`
	Price    float64 `json:"price" validate:"required"`
	Category string  `json:"category" validate:"required"`
}

var AllowedKeys = map[string]struct{}{
	"category": {},
}

func (q *ProductRequestQueryParams) ValidateQuery(c *fiber.Ctx) error {
	queryParams := c.Queries()

	for key := range queryParams {
		if _, ok := AllowedKeys[key]; !ok {
			return validator.ValidationErrors(nil)
		}
	}

	// check null query
	if _, hasCategoryKey := queryParams["category"]; hasCategoryKey {
		if q.Category == "" {
			return validator.ValidationErrors(nil)
		}
	}

	return nil
}
