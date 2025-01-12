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
type Repositories struct {
	Client         ClientRepository
	Product        ProductRepository
	Order          OrderRepository
	OrdersByClient OrdersByClientRepository
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
	Create(ctx context.Context, client *models.Client) error
	GetByID(ctx context.Context, id int64) (*models.Client, error)
	GetAll(ctx context.Context) ([]models.Client, error)
	Update(ctx context.Context, client *models.Client) error
	Delete(ctx context.Context, id int64) error
	GetClientOrders(ctx context.Context, id int64) ([]models.Order, error)
}

// ProductRepository определяет методы для работы с товарами
type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
	GetByID(ctx context.Context, id int64) (*models.Product, error)
	GetAll(ctx context.Context) ([]models.Product, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id int64) error
	GetProductOrderItems(ctx context.Context, id int64) ([]models.OrderItem, error)
}

// OrderRepository определяет методы для работы с заказами
type OrderRepository interface {
	Create(ctx context.Context, order *models.Order, items []models.OrderItem) error
	GetByID(ctx context.Context, id int64) (*models.Order, error)
	GetAll(ctx context.Context) ([]models.Order, error)
	Update(ctx context.Context, order *models.Order, items []models.OrderItem) error
	Delete(ctx context.Context, id int64) error
	Confirm(ctx context.Context, id int64) error
}

// OrdersByClientRepository определяет методы для работы с агрегированными суммами
type OrdersByClientRepository interface {
	GetByID(ctx context.Context, clientID int64) (*models.OrdersByClient, error)
	GetAll(ctx context.Context) ([]models.OrdersByClient, error)
	UpdateSum(ctx context.Context, clientID int64, amount float64) error
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
