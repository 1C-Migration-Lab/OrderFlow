package api

import (
	"github.com/1C-Migration-Lab/OrderFlow/internal/api/handlers"
	"github.com/1C-Migration-Lab/OrderFlow/internal/service"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, services *service.Services) {
	// Clients
	r.GET("/api/clients", handlers.GetClients(services.Client))
	r.GET("/api/clients/:id", handlers.GetClientByID(services.Client))
	r.POST("/api/clients", handlers.CreateClient(services.Client))
	r.PUT("/api/clients/:id", handlers.UpdateClient(services.Client))
	r.DELETE("/api/clients/:id", handlers.DeleteClient(services.Client))
	r.GET("/api/clients/:id/orders", handlers.GetClientOrders(services.Client))

	// Products
	r.GET("/api/products", handlers.GetProducts(services.Product))
	r.GET("/api/products/:id", handlers.GetProductByID(services.Product))
	r.POST("/api/products", handlers.CreateProduct(services.Product))
	r.PUT("/api/products/:id", handlers.UpdateProduct(services.Product))
	r.DELETE("/api/products/:id", handlers.DeleteProduct(services.Product))
	r.GET("/api/products/:id/order-items", handlers.GetProductOrderItems(services.Product))

	// Orders
	r.GET("/api/orders", handlers.GetOrders(services.Order))
	r.GET("/api/orders/:id", handlers.GetOrderByID(services.Order))
	r.POST("/api/orders", handlers.CreateOrder(services.Order))
	r.PUT("/api/orders/:id", handlers.UpdateOrder(services.Order))
	r.DELETE("/api/orders/:id", handlers.DeleteOrder(services.Order))
	r.POST("/api/orders/:id/confirm", handlers.ConfirmOrder(services.Order))

	// OrdersByClient
	r.GET("/api/orders-by-client", handlers.GetOrdersByClient(services.OrdersByClient))
	r.GET("/api/orders-by-client/:clientId", handlers.GetOrdersByClientID(services.OrdersByClient))
}
