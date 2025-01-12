package repository

import (
	"context"
	"database/sql"

	"github.com/1C-Migration-Lab/OrderFlow/internal/domain/models"
)

func (r *productRepository) Create(ctx context.Context, product *models.Product) error {
	query := `
		INSERT INTO products (name, unit)
		VALUES ($1, $2)
		RETURNING id`

	err := r.db.QueryRowContext(ctx, query, product.Name, product.Unit).Scan(&product.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *productRepository) GetByID(ctx context.Context, id int64) (*models.Product, error) {
	query := `
		SELECT id, name, unit
		FROM products
		WHERE id = $1`

	product := &models.Product{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&product.ID, &product.Name, &product.Unit)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *productRepository) Update(ctx context.Context, product *models.Product) error {
	query := `
		UPDATE products
		SET name = $1, unit = $2
		WHERE id = $3`

	result, err := r.db.ExecContext(ctx, query, product.Name, product.Unit, product.ID)
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

func (r *productRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM products WHERE id = $1`

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

func (r *productRepository) GetAll(ctx context.Context) ([]models.Product, error) {
	query := `
		SELECT id, name, unit
		FROM products
		ORDER BY name`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Unit); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// GetProductOrderItems возвращает все позиции заказов, где используется товар
func (r *productRepository) GetProductOrderItems(ctx context.Context, productID int64) ([]models.OrderItem, error) {
	query := `
		SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.price, oi.line_amount
		FROM order_items oi
		WHERE oi.product_id = $1`

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.Price,
			&item.LineAmount,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
