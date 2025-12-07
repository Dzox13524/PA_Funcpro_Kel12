package repository

import (
	"context"

	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/domain"
	"gorm.io/gorm"
)

type CreateProductRepoFunc func(ctx context.Context, p domain.Product) (domain.Product, error)
type GetAllProductsRepoFunc func(ctx context.Context) ([]domain.Product, error)
type GetProductByIDRepoFunc func(ctx context.Context, id string) (domain.Product, error)
type UpdateProductRepoFunc func(ctx context.Context, id string, updates map[string]interface{}) (domain.Product, error)
type DeleteProductRepoFunc func(ctx context.Context, id string) error

func NewCreateProductRepository(db *gorm.DB) CreateProductRepoFunc {
	return func(ctx context.Context, p domain.Product) (domain.Product, error) {
		result := db.WithContext(ctx).Create(&p)
		return p, result.Error
	}
}

func NewGetAllProductsRepository(db *gorm.DB) GetAllProductsRepoFunc {
	return func(ctx context.Context) ([]domain.Product, error) {
		var products []domain.Product
		result := db.WithContext(ctx).Find(&products)
		return products, result.Error
	}
}

func NewGetProductByIDRepository(db *gorm.DB) GetProductByIDRepoFunc {
	return func(ctx context.Context, id string) (domain.Product, error) {
		var p domain.Product
		result := db.WithContext(ctx).First(&p, "id = ?", id)
		return p, result.Error
	}
}

func NewUpdateProductRepository(db *gorm.DB) UpdateProductRepoFunc {
	return func(ctx context.Context, id string, updates map[string]interface{}) (domain.Product, error) {
		var p domain.Product
		if err := db.WithContext(ctx).First(&p, "id = ?", id).Error; err != nil {
			return domain.Product{}, err
		}

		if err := db.WithContext(ctx).Model(&p).Updates(updates).Error; err != nil {
			return domain.Product{}, err
		}
		return p, nil
	}
}

func NewDeleteProductRepository(db *gorm.DB) DeleteProductRepoFunc {
	return func(ctx context.Context, id string) error {
		return db.WithContext(ctx).Delete(&domain.Product{}, "id = ?", id).Error
	}
}