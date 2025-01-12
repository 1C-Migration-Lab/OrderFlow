package repository

import (
	"context"
	"database/sql"

	"github.com/1C-Migration-Lab/OrderFlow/internal/domain/models"
)

func (r *orderRepository) Create(ctx context.Context, order *models.Order, items []models.OrderItem) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Создаем заказ
	query := `
		INSERT INTO orders (client_id, date, number, total_amount, is_confirmed)
		VALUES ($1, CURRENT_TIMESTAMP, $2, 0, false)
		RETURNING id, date, created_at`

	err = tx.QueryRowContext(ctx, query, order.ClientID, order.Number).
		Scan(&order.ID, &order.Date, &order.CreatedAt)
	if err != nil {
		return err
	}

	// Создаем позиции заказа
	for i := range items {
		items[i].OrderID = order.ID
		query = `
			INSERT INTO order_items (order_id, product_id, quantity, price)
			VALUES ($1, $2, $3, $4)
			RETURNING id, line_amount`

		err = tx.QueryRowContext(ctx, query,
			items[i].OrderID,
			items[i].ProductID,
			items[i].Quantity,
			items[i].Price,
		).Scan(&items[i].ID, &items[i].LineAmount)
		if err != nil {
			return err
		}
	}

	// Обновляем общую сумму заказа
	query = `
		UPDATE orders
		SET total_amount = (
			SELECT COALESCE(SUM(line_amount), 0)
			FROM order_items
			WHERE order_id = $1
		)
		WHERE id = $1
		RETURNING total_amount`

	err = tx.QueryRowContext(ctx, query, order.ID).Scan(&order.TotalAmount)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *orderRepository) GetByID(ctx context.Context, id int64) (*models.Order, error) {
	// Получаем заказ
	query := `
		SELECT o.id, o.client_id, o.date, o.number, o.total_amount, o.is_confirmed, o.created_at,
			   c.id, c.name, c.inn
		FROM orders o
		JOIN clients c ON c.id = o.client_id
		WHERE o.id = $1`

	order := &models.Order{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&order.ID, &order.ClientID, &order.Date, &order.Number,
		&order.TotalAmount, &order.IsConfirmed, &order.CreatedAt,
		&order.Client.ID, &order.Client.Name, &order.Client.INN,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	// Получаем позиции заказа
	query = `
		SELECT i.id, i.product_id, i.quantity, i.price, i.line_amount,
			   p.id, p.name, p.unit
		FROM order_items i
		JOIN products p ON p.id = i.product_id
		WHERE i.order_id = $1`

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.OrderItem
		err := rows.Scan(
			&item.ID, &item.ProductID, &item.Quantity, &item.Price, &item.LineAmount,
			&item.Product.ID, &item.Product.Name, &item.Product.Unit,
		)
		if err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return order, nil
}

func (r *orderRepository) Update(ctx context.Context, order *models.Order, items []models.OrderItem) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Проверяем, что заказ не подтвержден
	var isConfirmed bool
	err = tx.QueryRowContext(ctx, "SELECT is_confirmed FROM orders WHERE id = $1", order.ID).
		Scan(&isConfirmed)
	if err == sql.ErrNoRows {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if isConfirmed {
		return ErrConflict
	}

	// Обновляем заказ
	query := `
		UPDATE orders
		SET client_id = $1, number = $2
		WHERE id = $3`

	result, err := tx.ExecContext(ctx, query, order.ClientID, order.Number, order.ID)
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

	// Удаляем старые позиции
	_, err = tx.ExecContext(ctx, "DELETE FROM order_items WHERE order_id = $1", order.ID)
	if err != nil {
		return err
	}

	// Создаем новые позиции
	for i := range items {
		items[i].OrderID = order.ID
		query = `
			INSERT INTO order_items (order_id, product_id, quantity, price)
			VALUES ($1, $2, $3, $4)
			RETURNING id, line_amount`

		err = tx.QueryRowContext(ctx, query,
			items[i].OrderID,
			items[i].ProductID,
			items[i].Quantity,
			items[i].Price,
		).Scan(&items[i].ID, &items[i].LineAmount)
		if err != nil {
			return err
		}
	}

	// Обновляем общую сумму заказа
	query = `
		UPDATE orders
		SET total_amount = (
			SELECT COALESCE(SUM(line_amount), 0)
			FROM order_items
			WHERE order_id = $1
		)
		WHERE id = $1
		RETURNING total_amount`

	err = tx.QueryRowContext(ctx, query, order.ID).Scan(&order.TotalAmount)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *orderRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM orders WHERE id = $1`

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

func (r *orderRepository) GetAll(ctx context.Context) ([]models.Order, error) {
	query := `
		SELECT o.id, o.client_id, o.date, o.number, o.total_amount, o.is_confirmed, o.created_at,
			   c.id, c.name, c.inn
		FROM orders o
		JOIN clients c ON c.id = o.client_id
		ORDER BY o.created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID, &order.ClientID, &order.Date, &order.Number,
			&order.TotalAmount, &order.IsConfirmed, &order.CreatedAt,
			&order.Client.ID, &order.Client.Name, &order.Client.INN,
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

func (r *orderRepository) Confirm(ctx context.Context, id int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Получаем данные заказа
	var clientID int64
	var totalAmount float64
	var isConfirmed bool

	query := `
		SELECT client_id, total_amount, is_confirmed
		FROM orders
		WHERE id = $1
		FOR UPDATE`

	err = tx.QueryRowContext(ctx, query, id).Scan(&clientID, &totalAmount, &isConfirmed)
	if err == sql.ErrNoRows {
		return ErrNotFound
	}
	if err != nil {
		return err
	}

	if isConfirmed {
		return ErrConflict
	}

	// Обновляем статус заказа
	query = `
		UPDATE orders
		SET is_confirmed = true
		WHERE id = $1`

	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	// Обновляем сумму в orders_by_client
	query = `
		INSERT INTO orders_by_client (client_id, orders_sum)
		VALUES ($1, $2)
		ON CONFLICT (client_id)
		DO UPDATE SET orders_sum = orders_by_client.orders_sum + $2`

	_, err = tx.ExecContext(ctx, query, clientID, totalAmount)
	if err != nil {
		return err
	}

	return tx.Commit()
}
