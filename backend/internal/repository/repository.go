package repository

import (
	"context"
	"errors"

	"github.com/your-project/backend/internal/domain/models"
)

var (
	ErrNotFound = errors.New("record not found")
	ErrConflict = errors.New("record already exists")
)

// Repository определяет интерфейс для работы с данными
type Repository interface {
	ClientRepository
	ProductRepository
	OrderRepository
	OrdersByClientRepository
}

// ClientRepository определяет методы для работы с клиентами
type ClientRepository interface {
	CreateClient(ctx context.Context, client *models.Client) error
	GetClient(ctx context.Context, id int64) (*models.Client, error)
	UpdateClient(ctx context.Context, client *models.Client) error
	DeleteClient(ctx context.Context, id int64) error
	ListClients(ctx context.Context) ([]models.Client, error)
}

// ProductRepository определяет методы для работы с товарами
type ProductRepository interface {
	CreateProduct(ctx context.Context, product *models.Product) error
	GetProduct(ctx context.Context, id int64) (*models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product) error
	DeleteProduct(ctx context.Context, id int64) error
	ListProducts(ctx context.Context) ([]models.Product, error)
}

// OrderRepository определяет методы для работы с заказами
type OrderRepository interface {
	CreateOrder(ctx context.Context, order *models.Order, items []models.OrderItem) error
	GetOrder(ctx context.Context, id int64) (*models.Order, error)
	UpdateOrder(ctx context.Context, order *models.Order, items []models.OrderItem) error
	DeleteOrder(ctx context.Context, id int64) error
	ListOrders(ctx context.Context) ([]models.Order, error)
	ConfirmOrder(ctx context.Context, id int64) error
}

// OrdersByClientRepository определяет методы для работы с агрегированными суммами
type OrdersByClientRepository interface {
	GetOrdersByClient(ctx context.Context, clientID int64) (*models.OrdersByClient, error)
	ListOrdersByClient(ctx context.Context) ([]models.OrdersByClient, error)
	UpdateOrdersSum(ctx context.Context, clientID int64, amount float64) error
}

// PostgresRepository реализует Repository для PostgreSQL
type PostgresRepository struct {
	// db *gorm.DB
}

// NewPostgresRepository создает новый экземпляр PostgresRepository
func NewPostgresRepository( /* db *gorm.DB */ ) Repository {
	// TODO: Реализовать все методы интерфейсов
	panic("not implemented")
}
