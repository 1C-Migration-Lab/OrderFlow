package repository

import (
	"context"
	"database/sql"

	"github.com/1C-Migration-Lab/OrderFlow/internal/domain/models"
)

func (r *clientRepository) Create(ctx context.Context, client *models.Client) error {
	query := `
		INSERT INTO clients (name, inn)
		VALUES ($1, $2)
		RETURNING id`

	err := r.db.QueryRowContext(ctx, query, client.Name, client.INN).Scan(&client.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *clientRepository) GetByID(ctx context.Context, id int64) (*models.Client, error) {
	query := `
		SELECT id, name, inn
		FROM clients
		WHERE id = $1`

	client := &models.Client{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&client.ID, &client.Name, &client.INN)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (r *clientRepository) Update(ctx context.Context, client *models.Client) error {
	query := `
		UPDATE clients
		SET name = $1, inn = $2
		WHERE id = $3`

	result, err := r.db.ExecContext(ctx, query, client.Name, client.INN, client.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *clientRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM clients WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *clientRepository) GetAll(ctx context.Context) ([]models.Client, error) {
	query := `
		SELECT id, name, inn
		FROM clients
		ORDER BY name`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []models.Client
	for rows.Next() {
		var client models.Client
		if err := rows.Scan(&client.ID, &client.Name, &client.INN); err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return clients, nil
}

// GetClientOrders возвращает все заказы клиента
func (r *clientRepository) GetClientOrders(ctx context.Context, clientID int64) ([]models.Order, error) {
	query := `
		SELECT o.id, o.client_id, o.date, o.number, o.total_amount, o.is_confirmed, o.created_at
		FROM orders o
		WHERE o.client_id = $1`

	rows, err := r.db.QueryContext(ctx, query, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID,
			&order.ClientID,
			&order.Date,
			&order.Number,
			&order.TotalAmount,
			&order.IsConfirmed,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
