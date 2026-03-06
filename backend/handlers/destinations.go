package handlers

import (
	"context"
	"net/http"
	"time"

	"backend/models"
	"backend/service"

	"github.com/gin-gonic/gin"
)

type DestinationHandler struct {
	destService *service.DestinationService
}

func NewDestinationHandler(destService *service.DestinationService) *DestinationHandler {
	return &DestinationHandler{destService: destService}
}

func (h *DestinationHandler) GetDestinations(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filters := service.ListFilters{
		Name:      c.Query("name"),
		Country:   c.Query("country"),
		MinRating: c.Query("minRating"),
		Page:      c.DefaultQuery("page", "1"),
		Limit:     c.DefaultQuery("limit", "10"),
	}

	result, err := h.destService.ListDestinations(ctx, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"destinations": result.Destinations,
		"pagination": gin.H{
			"total":      result.Pagination.Total,
			"page":       result.Pagination.Page,
			"limit":      result.Pagination.Limit,
			"totalPages": result.Pagination.TotalPages,
		},
	})
}

func (h *DestinationHandler) GetDestination(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dest, reviews, err := h.destService.GetDestination(ctx, c.Param("id"))
	if err != nil {
		switch err.Error() {
		case "invalid destination ID":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "destination not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"destination": dest,
		"reviews":     reviews,
	})
}

func (h *DestinationHandler) CreateDestination(c *gin.Context) {
	var input models.DestinationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dest, err := h.destService.CreateDestination(ctx, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"destination": dest})
}

func (h *DestinationHandler) UpdateDestination(c *gin.Context) {
	var input models.DestinationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := h.destService.UpdateDestination(ctx, c.Param("id"), input); err != nil {
		switch err.Error() {
		case "invalid destination ID":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "destination not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Destination updated successfully"})
}

func (h *DestinationHandler) DeleteDestination(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := h.destService.DeleteDestination(ctx, c.Param("id")); err != nil {
		switch err.Error() {
		case "invalid destination ID":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "destination not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Destination deleted successfully"})
}
