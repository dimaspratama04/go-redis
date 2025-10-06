package usecase

import (
	"golang-redis/internal/entity"
	"golang-redis/internal/repository"

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

func (uc *ProductUseCase) GetAllProducts() ([]entity.Product, error) {
	return uc.Repository.GetAll()
}
