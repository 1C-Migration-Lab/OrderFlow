package repository

import (
	"context"

	"github.com/1C-Migration-Lab/OrderFlow/internal/domain/models"
)

func (r *clientRepository) CreateClient(ctx context.Context, client *models.Client) error {
	// TODO: реализация
	return nil
}

func (r *clientRepository) GetClient(ctx context.Context, id int64) (*models.Client, error) {
	// TODO: реализация
	return nil, nil
}

func (r *clientRepository) UpdateClient(ctx context.Context, client *models.Client) error {
	// TODO: реализация
	return nil
}

func (r *clientRepository) DeleteClient(ctx context.Context, id int64) error {
	// TODO: реализация
	return nil
}

func (r *clientRepository) ListClients(ctx context.Context) ([]models.Client, error) {
	// TODO: реализация
	return nil, nil
}
