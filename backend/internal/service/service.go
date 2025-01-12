package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/1C-Migration-Lab/OrderFlow/internal/domain/models"
	"github.com/1C-Migration-Lab/OrderFlow/internal/repository"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrAlreadyConfirmed = errors.New("order already confirmed")
	ErrValidation       = errors.New("validation error")
	ErrOrderHasNoItems  = errors.New("order has no items")
	ErrInvalidQuantity  = errors.New("invalid quantity")
	ErrInvalidPrice     = errors.New("invalid price")
	ErrClientHasOrders  = errors.New("client has associated orders")
	ErrProductHasOrders = errors.New("product has associated orders")
	ErrDuplicateNumber  = errors.New("order number already exists")
	ErrConfirmedNoEdit  = errors.New("confirmed order cannot be edited")
)

type ClientService interface {
	Create(ctx context.Context, client *models.Client) error
	GetByID(ctx context.Context, id int64) (*models.Client, error)
	GetAll(ctx context.Context) ([]models.Client, error)
	Update(ctx context.Context, client *models.Client) error
	Delete(ctx context.Context, id int64) error
}

type ProductService interface {
	Create(ctx context.Context, product *models.Product) error
	GetByID(ctx context.Context, id int64) (*models.Product, error)
	GetAll(ctx context.Context) ([]models.Product, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id int64) error
}

type OrderService interface {
	Create(ctx context.Context, order *models.Order, items []models.OrderItem) error
	GetByID(ctx context.Context, id int64) (*models.Order, error)
	GetAll(ctx context.Context) ([]models.Order, error)
	Update(ctx context.Context, order *models.Order, items []models.OrderItem) error
	Delete(ctx context.Context, id int64) error
	Confirm(ctx context.Context, id int64) error
}

type OrdersByClientService interface {
	GetByID(ctx context.Context, clientID int64) (*models.OrdersByClient, error)
	GetAll(ctx context.Context) ([]models.OrdersByClient, error)
}

type Services struct {
	Client         ClientService
	Product        ProductService
	Order          OrderService
	OrdersByClient OrdersByClientService
}

func NewServices(repos *repository.PostgresRepositories) *Services {
	return &Services{
		Client:         NewClientService(repos.Client),
		Product:        NewProductService(repos.Product),
		Order:          NewOrderService(repos.Order, repos.OrdersByClient),
		OrdersByClient: NewOrdersByClientService(repos.OrdersByClient),
	}
}

// ClientService implementation
type clientService struct {
	repo repository.ClientRepository
}

func NewClientService(repo repository.ClientRepository) ClientService {
	return &clientService{repo: repo}
}

func (s *clientService) Create(ctx context.Context, client *models.Client) error {
	if client.Name == "" {
		return fmt.Errorf("%w: name is required", ErrValidation)
	}
	return s.repo.Create(ctx, client)
}

func (s *clientService) GetByID(ctx context.Context, id int64) (*models.Client, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *clientService) GetAll(ctx context.Context) ([]models.Client, error) {
	return s.repo.GetAll(ctx)
}

func (s *clientService) Update(ctx context.Context, client *models.Client) error {
	if client.Name == "" {
		return fmt.Errorf("%w: name is required", ErrValidation)
	}
	return s.repo.Update(ctx, client)
}

func (s *clientService) Delete(ctx context.Context, id int64) error {
	// Сначала получаем клиента, чтобы проверить его существование
	client, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if client == nil {
		return ErrNotFound
	}
	return s.repo.Delete(ctx, id)
}

// ProductService implementation
type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) Create(ctx context.Context, product *models.Product) error {
	if product.Name == "" || product.Unit == "" {
		return fmt.Errorf("%w: name and unit are required", ErrValidation)
	}
	return s.repo.Create(ctx, product)
}

func (s *productService) GetByID(ctx context.Context, id int64) (*models.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *productService) GetAll(ctx context.Context) ([]models.Product, error) {
	return s.repo.GetAll(ctx)
}

func (s *productService) Update(ctx context.Context, product *models.Product) error {
	if product.Name == "" || product.Unit == "" {
		return fmt.Errorf("%w: name and unit are required", ErrValidation)
	}
	return s.repo.Update(ctx, product)
}

func (s *productService) Delete(ctx context.Context, id int64) error {
	// Сначала получаем продукт, чтобы проверить его существование
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if product == nil {
		return ErrNotFound
	}
	return s.repo.Delete(ctx, id)
}

// OrderService implementation
type orderService struct {
	repo         repository.OrderRepository
	ordersByRepo repository.OrdersByClientRepository
}

func NewOrderService(repo repository.OrderRepository, ordersByRepo repository.OrdersByClientRepository) OrderService {
	return &orderService{
		repo:         repo,
		ordersByRepo: ordersByRepo,
	}
}

func (s *orderService) Create(ctx context.Context, order *models.Order, items []models.OrderItem) error {
	if len(items) == 0 {
		return ErrOrderHasNoItems
	}

	// Валидация
	for _, item := range items {
		if item.Quantity <= 0 {
			return ErrInvalidQuantity
		}
		if item.Price <= 0 {
			return ErrInvalidPrice
		}
	}

	return s.repo.Create(ctx, order, items)
}

func (s *orderService) GetByID(ctx context.Context, id int64) (*models.Order, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *orderService) GetAll(ctx context.Context) ([]models.Order, error) {
	return s.repo.GetAll(ctx)
}

func (s *orderService) Update(ctx context.Context, order *models.Order, items []models.OrderItem) error {
	if order.IsConfirmed {
		return ErrConfirmedNoEdit
	}

	if len(items) > 0 {
		for _, item := range items {
			if item.Quantity <= 0 {
				return ErrInvalidQuantity
			}
			if item.Price <= 0 {
				return ErrInvalidPrice
			}
		}
		order.TotalAmount = 0
		for _, item := range items {
			order.TotalAmount += item.Quantity * item.Price
		}
	}

	return s.repo.Update(ctx, order, items)
}

func (s *orderService) Delete(ctx context.Context, id int64) error {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Если заказ подтвержден, корректируем суммы в OrdersByClient
	if order.IsConfirmed {
		if err := s.ordersByRepo.UpdateSum(ctx, order.ClientID, -order.TotalAmount); err != nil {
			return err
		}
	}

	return s.repo.Delete(ctx, id)
}

func (s *orderService) Confirm(ctx context.Context, id int64) error {
	order, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if order.IsConfirmed {
		return ErrAlreadyConfirmed
	}

	if len(order.Items) == 0 {
		return ErrOrderHasNoItems
	}

	// Обновляем сумму заказов клиента
	if err := s.ordersByRepo.UpdateSum(ctx, order.ClientID, order.TotalAmount); err != nil {
		return err
	}

	return s.repo.Confirm(ctx, id)
}

// OrdersByClientService implementation
type ordersByClientService struct {
	repo repository.OrdersByClientRepository
}

func NewOrdersByClientService(repo repository.OrdersByClientRepository) OrdersByClientService {
	return &ordersByClientService{repo: repo}
}

func (s *ordersByClientService) GetByID(ctx context.Context, clientID int64) (*models.OrdersByClient, error) {
	return s.repo.GetByID(ctx, clientID)
}

func (s *ordersByClientService) GetAll(ctx context.Context) ([]models.OrdersByClient, error) {
	return s.repo.GetAll(ctx)
}
