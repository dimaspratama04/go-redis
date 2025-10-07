package repository

import (
	"context"
	"golang-redis/internal/entity"

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

func (p *ProductRepository) CreateBatch(products []entity.Product) error {
	return p.DB.Create(&products).Error
}

func (p *ProductRepository) GetAll() ([]entity.Product, error) {
	var products []entity.Product
	if err := p.DB.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (p *ProductRepository) GetByID(id int) (*entity.Product, error) {
	var product entity.Product
	if err := p.DB.Where("id = ?", id).First(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}
