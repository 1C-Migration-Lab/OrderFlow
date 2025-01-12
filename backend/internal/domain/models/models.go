package models

import (
	"time"
)

// Client представляет клиента в системе
type Client struct {
	ID   int64  `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"not null"`
	INN  string `json:"inn"`
}

// Product представляет товар в системе
type Product struct {
	ID   int64  `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"not null"`
	Unit string `json:"unit" gorm:"not null"`
}

// Order представляет заказ в системе
type Order struct {
	ID          int64       `json:"id" gorm:"primaryKey"`
	ClientID    int64       `json:"client_id" gorm:"not null"`
	Client      Client      `json:"client" gorm:"foreignKey:ClientID"`
	Date        time.Time   `json:"date" gorm:"not null;default:CURRENT_TIMESTAMP"`
	Number      string      `json:"number" gorm:"not null;unique"`
	TotalAmount float64     `json:"total_amount" gorm:"type:decimal(15,2);not null;default:0"`
	IsConfirmed bool        `json:"is_confirmed" gorm:"not null;default:false"`
	CreatedAt   time.Time   `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
	Items       []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
}

// OrderItem представляет позицию заказа
type OrderItem struct {
	ID         int64   `json:"id" gorm:"primaryKey"`
	OrderID    int64   `json:"order_id" gorm:"not null"`
	ProductID  int64   `json:"product_id" gorm:"not null"`
	Product    Product `json:"product" gorm:"foreignKey:ProductID"`
	Quantity   float64 `json:"quantity" gorm:"type:decimal(15,3);not null;check:quantity >= 0"`
	Price      float64 `json:"price" gorm:"type:decimal(15,2);not null;check:price >= 0"`
	LineAmount float64 `json:"line_amount" gorm:"type:decimal(15,2);not null;default:0"`
}

// OrdersByClient представляет агрегированные суммы заказов по клиентам
type OrdersByClient struct {
	ClientID  int64   `json:"client_id" gorm:"primaryKey"`
	Client    Client  `json:"client" gorm:"foreignKey:ClientID"`
	OrdersSum float64 `json:"orders_sum" gorm:"type:decimal(15,2);not null;default:0"`
}

// CreateOrderRequest представляет запрос на создание заказа
type CreateOrderRequest struct {
	ClientID int64             `json:"client_id" binding:"required"`
	Number   string            `json:"number" binding:"required"`
	Items    []CreateOrderItem `json:"items" binding:"required,dive"`
}

// CreateOrderItem представляет запрос на создание позиции заказа
type CreateOrderItem struct {
	ProductID int64   `json:"product_id" binding:"required"`
	Quantity  float64 `json:"quantity" binding:"required,min=0"`
	Price     float64 `json:"price" binding:"required,min=0"`
}

// UpdateOrderRequest представляет запрос на обновление заказа
type UpdateOrderRequest struct {
	ClientID int64             `json:"client_id"`
	Number   string            `json:"number"`
	Items    []CreateOrderItem `json:"items" binding:"dive"`
}
