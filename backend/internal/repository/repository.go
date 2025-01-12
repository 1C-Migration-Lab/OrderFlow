package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/1C-Migration-Lab/OrderFlow/internal/domain/models"
)

var (
	ErrNotFound = errors.New("record not found")
	ErrConflict = errors.New("record already exists")
)

// Repository определяет интерфейс для работы с данными
type Repositories interface {
	ClientRepository
	ProductRepository
	OrderRepository
	OrdersByClientRepository
}

type PostgresRepositories struct {
	Client         ClientRepository
	Product        ProductRepository
	Order          OrderRepository
	OrdersByClient OrdersByClientRepository
}

func NewPostgresRepository(db *sql.DB) *PostgresRepositories {
	return &PostgresRepositories{
		Client:         NewClientRepository(db),
		Product:        NewProductRepository(db),
		Order:          NewOrderRepository(db),
		OrdersByClient: NewOrdersByClientRepository(db),
	}
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

// Структуры конкретных репозиториев
type clientRepository struct {
	db *sql.DB
}

type productRepository struct {
	db *sql.DB
}

type orderRepository struct {
	db *sql.DB
}

type ordersByClientRepository struct {
	db *sql.DB
}

// Функции создания репозиториев
func NewClientRepository(db *sql.DB) ClientRepository {
	return &clientRepository{
		db: db,
	}
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{
		db: db,
	}
}

func NewOrdersByClientRepository(db *sql.DB) OrdersByClientRepository {
	return &ordersByClientRepository{
		db: db,
	}
}
