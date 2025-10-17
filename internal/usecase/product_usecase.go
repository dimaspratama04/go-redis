package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"golang-redis/internal/delivery/http/request"
	"golang-redis/internal/entity"
	"golang-redis/internal/repository"
	"time"
)

type ProductUseCase struct {
	Repository *repository.ProductRepository
}

func NewProductUseCase(repository *repository.ProductRepository) *ProductUseCase {
	return &ProductUseCase{Repository: repository}
}

func (uc *ProductUseCase) CreateProduct(name string, price float64) (*entity.Product, error) {
	return uc.Repository.Create(name, price)
}

func (uc *ProductUseCase) CreateProductBatch(productPayload []request.ProductRequest) ([]entity.Product, error) {
	products := make([]entity.Product, 0, len(productPayload))
	for _, p := range productPayload {
		products = append(products, entity.Product{
			Name:  p.Name,
			Price: p.Price,
		})
	}

	if err := uc.Repository.CreateBatch(products); err != nil {
		return nil, err
	}

	return products, nil
}

func (uc *ProductUseCase) GetAllProducts() ([]entity.Product, error) {
	ctx := context.Background()
	cacheKey := "products:all"

	// check cache from redis
	cached, err := uc.Repository.RDB.Get(ctx, cacheKey).Result()
	if err == nil && cached != "" {
		var products []entity.Product
		if err := json.Unmarshal([]byte(cached), &products); err == nil {
			return products, nil
		}
	}

	// get data from db if cache not exist
	products, err := uc.Repository.GetAll()
	if err != nil {
		return nil, err
	}

	// store to redis 1 hour
	jsonData, _ := json.Marshal(products)
	uc.Repository.RDB.Set(ctx, cacheKey, jsonData, 1*time.Hour)

	return products, nil
}

func (uc *ProductUseCase) GetProductByID(ctx context.Context, id int) (*entity.Product, error) {
	cacheKeyProductId := fmt.Sprintf("product:%d", id)
	val, err := uc.Repository.RDB.Get(ctx, cacheKeyProductId).Result()
	if err == nil {
		var product entity.Product
		if err := json.Unmarshal([]byte(val), &product); err == nil {
			return &product, nil
		}
	}

	// if cache doesnt exist, get from db
	product, err := uc.Repository.GetByID(id)
	if err != nil {
		return nil, err
	}

	// store cache to redis 60s
	data, _ := json.Marshal(product)
	uc.Repository.RDB.Set(ctx, cacheKeyProductId, data, 1*time.Hour)
	return product, nil

}
