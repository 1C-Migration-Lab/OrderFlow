package repository

import (
	"context"

	"github.com/1C-Migration-Lab/OrderFlow/internal/domain/models"
)

func (r *orderRepository) CreateOrder(ctx context.Context, order *models.Order, items []models.OrderItem) error {
	// TODO: реализация
	return nil
}

func (r *orderRepository) GetOrder(ctx context.Context, id int64) (*models.Order, error) {
	// TODO: реализация
	return nil, nil
}

func (r *orderRepository) UpdateOrder(ctx context.Context, order *models.Order, items []models.OrderItem) error {
	// TODO: реализация
	return nil
}

func (r *orderRepository) DeleteOrder(ctx context.Context, id int64) error {
	// TODO: реализация
	return nil
}

func (r *orderRepository) ListOrders(ctx context.Context) ([]models.Order, error) {
	// TODO: реализация
	return nil, nil
}

func (r *orderRepository) ConfirmOrder(ctx context.Context, id int64) error {
	// TODO: реализация
	return nil
}
