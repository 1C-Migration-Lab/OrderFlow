package repository

import (
	"context"

	"github.com/1C-Migration-Lab/OrderFlow/internal/domain/models"
)

func (r *productRepository) CreateProduct(ctx context.Context, product *models.Product) error {
	// TODO: реализация
	return nil
}

func (r *productRepository) GetProduct(ctx context.Context, id int64) (*models.Product, error) {
	// TODO: реализация
	return nil, nil
}

func (r *productRepository) UpdateProduct(ctx context.Context, product *models.Product) error {
	// TODO: реализация
	return nil
}

func (r *productRepository) DeleteProduct(ctx context.Context, id int64) error {
	// TODO: реализация
	return nil
}

func (r *productRepository) ListProducts(ctx context.Context) ([]models.Product, error) {
	// TODO: реализация
	return nil, nil
}
