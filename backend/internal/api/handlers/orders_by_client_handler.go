package handlers

import (
	"net/http"
	"strconv"

	"github.com/1C-Migration-Lab/OrderFlow/internal/service"
	"github.com/gin-gonic/gin"
)

func GetOrdersByClient(s service.OrdersByClientService) gin.HandlerFunc {
	return func(c *gin.Context) {
		orders, err := s.GetAll(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, orders)
	}
}

func GetOrdersByClientID(s service.OrdersByClientService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("clientId"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
			return
		}

		orders, err := s.GetByID(c.Request.Context(), id)
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, orders)
	}
}
