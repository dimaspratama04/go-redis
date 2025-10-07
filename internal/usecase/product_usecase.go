package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"golang-redis/internal/delivery/http/request"
	"golang-redis/internal/entity"
	"golang-redis/internal/repository"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ProductUseCase struct {
	Repository *repository.ProductRepository
	DB         *gorm.DB
	RDB        *redis.Client
}

func NewProductUseCase(db *gorm.DB, rdb *redis.Client) *ProductUseCase {
	productRepo := repository.NewProductRepository(db, rdb)
	return &ProductUseCase{Repository: productRepo, DB: db, RDB: rdb}
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

func (uc *ProductUseCase) GetAllProducts(ctx context.Context) ([]entity.Product, error) {
	cacheKey := "products:all"

	// 1️⃣ Cek di Redis dulu
	cached, err := uc.RDB.Get(ctx, cacheKey).Result()
	if err == nil && cached != "" {
		var products []entity.Product
		if err := json.Unmarshal([]byte(cached), &products); err == nil {
			return products, nil
		}
	}

	// 2️⃣ Kalau cache miss → ambil dari DB
	products, err := uc.Repository.GetAll()
	if err != nil {
		return nil, err
	}

	// 3️⃣ Simpan hasil ke Redis 1 jam
	jsonData, _ := json.Marshal(products)
	uc.RDB.Set(ctx, cacheKey, jsonData, 1*time.Hour)

	return products, nil
}

func (uc *ProductUseCase) GetProductByID(ctx context.Context, id int) (*entity.Product, error) {
	cacheKeyProductId := fmt.Sprintf("product:%d", id)
	val, err := uc.RDB.Get(ctx, cacheKeyProductId).Result()
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
	uc.RDB.Set(ctx, cacheKeyProductId, data, 1*time.Hour)
	return product, nil

}
