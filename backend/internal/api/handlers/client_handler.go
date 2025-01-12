package handlers

import (
	"net/http"
	"strconv"

	"github.com/1C-Migration-Lab/OrderFlow/internal/domain/models"
	"github.com/1C-Migration-Lab/OrderFlow/internal/service"
	"github.com/gin-gonic/gin"
)

func GetClients(s service.ClientService) gin.HandlerFunc {
	return func(c *gin.Context) {
		clients, err := s.GetAll(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, clients)
	}
}

func GetClientByID(s service.ClientService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
			return
		}

		client, err := s.GetByID(c.Request.Context(), id)
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, client)
	}
}

func CreateClient(s service.ClientService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var client models.Client
		if err := c.ShouldBindJSON(&client); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := s.Create(c.Request.Context(), &client); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, client)
	}
}

func UpdateClient(s service.ClientService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
			return
		}

		var client models.Client
		if err := c.ShouldBindJSON(&client); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		client.ID = id

		if err := s.Update(c.Request.Context(), &client); err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, client)
	}
}

func DeleteClient(s service.ClientService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
			return
		}

		if err := s.Delete(c.Request.Context(), id); err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
			return
		} else if err == service.ErrClientHasOrders {
			c.JSON(http.StatusConflict, gin.H{"error": "client has associated orders"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

func GetClientOrders(s service.ClientService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
			return
		}

		orders, err := s.GetClientOrders(c.Request.Context(), id)
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
