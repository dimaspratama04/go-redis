package repository

import (
	"context"
	"encoding/json"
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

func (p *ProductRepository) Create(name string, price float64) (*entity.Product, error) {
	product := &entity.Product{Name: name, Price: price}
	if err := p.DB.Create(product).Error; err != nil {
		return nil, err
	}

	// delete cache if already insert new products
	_ = p.RDB.Del(context.Background(), "products:all").Err()
	return product, nil
}

func (p *ProductRepository) GetAll() ([]entity.Product, error) {
	var products []entity.Product
	cacheRedisKey := "products:all"

	// check cache first
	cached, err := p.RDB.Get(context.Background(), cacheRedisKey).Result()
	if err == nil && cached != "" {
		_ = json.Unmarshal([]byte(cached), &products)
		return products, nil
	}

	// check from db if cache miss
	if err := p.DB.Find(&products).Error; err != nil {
		return nil, err
	}

	// store to redis 60s
	productsJSON, _ := json.Marshal(products)
	p.RDB.Set(context.Background(), cacheRedisKey, productsJSON, 60*time.Second)

	return products, nil
}
