package handlers

import (
	"net/http"
	"strconv"

	"github.com/1C-Migration-Lab/OrderFlow/internal/domain/models"
	"github.com/1C-Migration-Lab/OrderFlow/internal/service"
	"github.com/gin-gonic/gin"
)

func GetProducts(s service.ProductService) gin.HandlerFunc {
	return func(c *gin.Context) {
		products, err := s.GetAll(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, products)
	}
}

func GetProductByID(s service.ProductService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
			return
		}

		product, err := s.GetByID(c.Request.Context(), id)
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, product)
	}
}

func CreateProduct(s service.ProductService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product models.Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := s.Create(c.Request.Context(), &product); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, product)
	}
}

func UpdateProduct(s service.ProductService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
			return
		}

		var product models.Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		product.ID = id

		if err := s.Update(c.Request.Context(), &product); err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, product)
	}
}

func DeleteProduct(s service.ProductService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
			return
		}

		if err := s.Delete(c.Request.Context(), id); err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		} else if err == service.ErrProductHasOrders {
			c.JSON(http.StatusConflict, gin.H{"error": "product has associated orders"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

func GetProductOrderItems(s service.ProductService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
			return
		}

		items, err := s.GetProductOrderItems(c.Request.Context(), id)
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, items)
	}
}
