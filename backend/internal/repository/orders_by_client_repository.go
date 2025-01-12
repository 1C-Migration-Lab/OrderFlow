package repository

import (
	"context"

	"github.com/1C-Migration-Lab/OrderFlow/internal/domain/models"
)

func (r *ordersByClientRepository) GetOrdersByClient(ctx context.Context, clientID int64) (*models.OrdersByClient, error) {
	// TODO: реализация
	return nil, nil
}

func (r *ordersByClientRepository) ListOrdersByClient(ctx context.Context) ([]models.OrdersByClient, error) {
	// TODO: реализация
	return nil, nil
}

func (r *ordersByClientRepository) UpdateOrdersSum(ctx context.Context, clientID int64, amount float64) error {
	// TODO: реализация
	return nil
}
