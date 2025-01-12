package repository

import (
	"context"
	"database/sql"

	"github.com/1C-Migration-Lab/OrderFlow/internal/domain/models"
)

func (r *ordersByClientRepository) GetByID(ctx context.Context, clientID int64) (*models.OrdersByClient, error) {
	query := `
		SELECT obc.client_id, obc.orders_sum,
			   c.id, c.name, c.inn
		FROM orders_by_client obc
		JOIN clients c ON c.id = obc.client_id
		WHERE obc.client_id = $1`

	result := &models.OrdersByClient{}
	err := r.db.QueryRowContext(ctx, query, clientID).Scan(
		&result.ClientID, &result.OrdersSum,
		&result.Client.ID, &result.Client.Name, &result.Client.INN,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *ordersByClientRepository) GetAll(ctx context.Context) ([]models.OrdersByClient, error) {
	query := `
		SELECT obc.client_id, obc.orders_sum,
			   c.id, c.name, c.inn
		FROM orders_by_client obc
		JOIN clients c ON c.id = obc.client_id
		ORDER BY obc.orders_sum DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.OrdersByClient
	for rows.Next() {
		var result models.OrdersByClient
		err := rows.Scan(
			&result.ClientID, &result.OrdersSum,
			&result.Client.ID, &result.Client.Name, &result.Client.INN,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *ordersByClientRepository) UpdateSum(ctx context.Context, clientID int64, amount float64) error {
	query := `
		INSERT INTO orders_by_client (client_id, orders_sum)
		VALUES ($1, $2)
		ON CONFLICT (client_id)
		DO UPDATE SET orders_sum = orders_by_client.orders_sum + $2`

	result, err := r.db.ExecContext(ctx, query, clientID, amount)
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
