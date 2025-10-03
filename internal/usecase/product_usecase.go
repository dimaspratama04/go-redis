package usecase

import (
	"golang-redis/internal/entity"

	"gorm.io/gorm"
)

type ProductUseCase struct {
	DB *gorm.DB
}

func NewProductUseCase(db *gorm.DB) *ProductUseCase {
	return &ProductUseCase{DB: db}
}

func (uc *ProductUseCase) CreateProduct(name string, price float64) (*entity.Product, error) {
	product := &entity.Product{Name: name, Price: price}
	if err := uc.DB.Create(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (uc *ProductUseCase) GetAllProducts() ([]entity.Product, error) {
	var products []entity.Product
	if err := uc.DB.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
