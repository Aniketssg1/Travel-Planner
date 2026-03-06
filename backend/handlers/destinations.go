package handlers

import (
	"net/http"

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
	filters := service.ListFilters{
		Name:      c.Query("name"),
		Country:   c.Query("country"),
		MinRating: c.Query("minRating"),
		Page:      c.DefaultQuery("page", "1"),
		Limit:     c.DefaultQuery("limit", "10"),
	}

	result, err := h.destService.ListDestinations(c.Request.Context(), filters)
	if err != nil {
		handleError(c, err)
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
	dest, reviews, err := h.destService.GetDestination(c.Request.Context(), c.Param("id"))
	if err != nil {
		handleError(c, err)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": "VALIDATION_ERROR"})
		return
	}

	dest, err := h.destService.CreateDestination(c.Request.Context(), input)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"destination": dest})
}

func (h *DestinationHandler) UpdateDestination(c *gin.Context) {
	var input models.DestinationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": "VALIDATION_ERROR"})
		return
	}

	if err := h.destService.UpdateDestination(c.Request.Context(), c.Param("id"), input); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Destination updated successfully"})
}

func (h *DestinationHandler) DeleteDestination(c *gin.Context) {
	if err := h.destService.DeleteDestination(c.Request.Context(), c.Param("id")); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Destination deleted successfully"})
}
