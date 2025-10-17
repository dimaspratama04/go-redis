package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"golang-redis/internal/entity"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ProductRepository struct {
	DB  *gorm.DB
	RDB *redis.Client
}

func NewProductRepository(db *gorm.DB, rdb *redis.Client) *ProductRepository {
	return &ProductRepository{DB: db, RDB: rdb}
}

func (p *ProductRepository) Create(ctx context.Context, name string, price float64, category string) (*entity.Product, error) {
	product := &entity.Product{Name: name, Price: price, Category: category}
	if err := p.DB.Create(product).Error; err != nil {
		return nil, err
	}

	// delete cache if already insert new product
	_ = p.RDB.Del(ctx, "products:all").Err()

	// delete cache by category
	cacheKeyByCategory := "products:type:" + category
	_ = p.RDB.Del(ctx, cacheKeyByCategory).Err()

	return product, nil
}

func (p *ProductRepository) CreateBatch(ctx context.Context, products []entity.Product) error {
	// delete cache if already insert single product
	_ = p.RDB.Del(ctx, "products:all").Err()

	// delete cache by category
	uniqueCategories := make(map[string]struct{})
	for _, product := range products {
		uniqueCategories[product.Category] = struct{}{}
	}

	var keysToDelete []string
	for category := range uniqueCategories {
		cacheKey := "products:type:" + category
		keysToDelete = append(keysToDelete, cacheKey)
	}

	if len(keysToDelete) > 0 {
		_ = p.RDB.Del(ctx, keysToDelete...).Err()
	}

	return p.DB.Create(&products).Error
}

func (p *ProductRepository) GetAll(ctx context.Context) ([]entity.Product, error) {
	var products []entity.Product
	cacheKey := "products:all"

	// check cache from redis
	cached, err := p.RDB.Get(ctx, cacheKey).Result()
	if err == nil && cached != "" {
		if err := json.Unmarshal([]byte(cached), &products); err == nil {
			return products, nil
		}
	}

	// check to database
	if err := p.DB.Find(&products).Error; err != nil {
		return nil, err
	}

	// store to redis 1 hour
	jsonData, _ := json.Marshal(products)
	p.RDB.Set(ctx, cacheKey, jsonData, 1*time.Hour)

	return products, nil
}

func (p *ProductRepository) GetByID(ctx context.Context, id int) (*entity.Product, error) {
	var product entity.Product

	// check to redis
	cacheKeyProductId := fmt.Sprintf("product:%d", id)
	val, err := p.RDB.Get(ctx, cacheKeyProductId).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(val), &product); err == nil {
			return &product, nil
		}
	}

	// check to database
	if err := p.DB.Where("id = ?", id).First(&product).Error; err != nil {
		return nil, err
	}

	// store to redis
	data, _ := json.Marshal(product)
	p.RDB.Set(ctx, cacheKeyProductId, data, 1*time.Hour)
	return &product, nil
}

func (p *ProductRepository) GetByCategory(ctx context.Context, productType string) ([]entity.Product, error) {
	var products []entity.Product
	cacheKey := fmt.Sprintf("products:type:%s", productType)

	// cek cache redis
	cached, err := p.RDB.Get(ctx, cacheKey).Result()
	if err == nil && cached != "" {
		if err := json.Unmarshal([]byte(cached), &products); err == nil {
			return products, nil
		}
	}

	// get from db
	if err := p.DB.Where("category = ?", productType).Find(&products).Error; err != nil {
		return nil, err
	}

	// simpan ke redis selama 1 jam
	jsonData, _ := json.Marshal(products)
	p.RDB.Set(ctx, cacheKey, jsonData, 1*time.Hour)

	return products, nil
}
