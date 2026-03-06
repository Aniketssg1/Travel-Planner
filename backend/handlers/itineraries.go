package handlers

import (
	"context"
	"net/http"
	"time"

	"backend/models"
	"backend/service"

	"github.com/gin-gonic/gin"
)

type ItineraryHandler struct {
	itinService *service.ItineraryService
}

func NewItineraryHandler(itinService *service.ItineraryService) *ItineraryHandler {
	return &ItineraryHandler{itinService: itinService}
}

func (h *ItineraryHandler) GetItineraries(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	itineraries, err := h.itinService.GetItineraries(ctx, c.GetString("userId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"itineraries": itineraries})
}

func (h *ItineraryHandler) CreateItinerary(c *gin.Context) {
	var input models.ItineraryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	itin, err := h.itinService.CreateItinerary(ctx, c.GetString("userId"), input)
	if err != nil {
		switch err.Error() {
		case "invalid startDate format, use YYYY-MM-DD",
			"invalid endDate format, use YYYY-MM-DD",
			"endDate must be after startDate":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"itinerary": itin})
}

func (h *ItineraryHandler) UpdateItinerary(c *gin.Context) {
	var input models.ItineraryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := h.itinService.UpdateItinerary(ctx, c.Param("id"), c.GetString("userId"), input); err != nil {
		switch err.Error() {
		case "invalid itinerary ID",
			"invalid startDate format, use YYYY-MM-DD",
			"invalid endDate format, use YYYY-MM-DD",
			"endDate must be after startDate":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "itinerary not found or access denied":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Itinerary updated successfully"})
}

func (h *ItineraryHandler) DeleteItinerary(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := h.itinService.DeleteItinerary(ctx, c.Param("id"), c.GetString("userId")); err != nil {
		switch err.Error() {
		case "invalid itinerary ID":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "itinerary not found or access denied":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Itinerary deleted successfully"})
}
