package handlers

import (
	"net/http"

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
	itineraries, err := h.itinService.GetItineraries(c.Request.Context(), c.GetString("userId"))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"itineraries": itineraries})
}

func (h *ItineraryHandler) CreateItinerary(c *gin.Context) {
	var input models.ItineraryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": "VALIDATION_ERROR"})
		return
	}

	itin, err := h.itinService.CreateItinerary(c.Request.Context(), c.GetString("userId"), input)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"itinerary": itin})
}

func (h *ItineraryHandler) UpdateItinerary(c *gin.Context) {
	var input models.ItineraryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": "VALIDATION_ERROR"})
		return
	}

	if err := h.itinService.UpdateItinerary(c.Request.Context(), c.Param("id"), c.GetString("userId"), input); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Itinerary updated successfully"})
}

func (h *ItineraryHandler) DeleteItinerary(c *gin.Context) {
	if err := h.itinService.DeleteItinerary(c.Request.Context(), c.Param("id"), c.GetString("userId")); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Itinerary deleted successfully"})
}
