package usecase

import (
	"context"
	"golang-redis/internal/delivery/http/request"
	"golang-redis/internal/entity"
	"golang-redis/internal/repository"
)

type ProductUseCase struct {
	Repository *repository.ProductRepository
}

func NewProductUseCase(repository *repository.ProductRepository) *ProductUseCase {
	return &ProductUseCase{Repository: repository}
}

func (uc *ProductUseCase) CreateProduct(ctx context.Context, name string, price float64, category string) (*entity.Product, error) {
	return uc.Repository.Create(ctx, name, price, category)
}

func (uc *ProductUseCase) CreateProductBatch(ctx context.Context, productPayload []request.ProductRequest) ([]entity.Product, error) {
	products := make([]entity.Product, 0, len(productPayload))
	for _, p := range productPayload {
		products = append(products, entity.Product{
			Name:     p.Name,
			Price:    p.Price,
			Category: p.Category,
		})
	}

	if err := uc.Repository.CreateBatch(ctx, products); err != nil {
		return nil, err
	}

	return products, nil
}

func (uc *ProductUseCase) GetAllProducts(ctx context.Context) ([]entity.Product, error) {
	products, err := uc.Repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (uc *ProductUseCase) GetProductByID(ctx context.Context, id int) (*entity.Product, error) {
	product, err := uc.Repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return product, nil

}

func (uc *ProductUseCase) GetProductByCategory(ctx context.Context, productCategory string) ([]entity.Product, error) {
	products, err := uc.Repository.GetByCategory(ctx, productCategory)
	if err != nil {
		return nil, err
	}

	return products, nil

}
