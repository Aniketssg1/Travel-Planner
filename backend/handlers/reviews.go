package handlers

import (
	"context"
	"net/http"
	"time"

	"backend/models"
	"backend/repository"
	"backend/service"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	reviewService *service.ReviewService
	userRepo      *repository.UserRepository
}

func NewReviewHandler(reviewService *service.ReviewService, userRepo *repository.UserRepository) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService, userRepo: userRepo}
}

func (h *ReviewHandler) SubmitReview(c *gin.Context) {
	var input models.ReviewInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userIDStr := c.GetString("userId")

	userName := "User"
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err == nil {
		user, err := h.userRepo.FindByID(ctx, userID)
		if err == nil {
			userName = user.Name
		}
	}

	review, err := h.reviewService.SubmitReviewWithUserName(ctx, c.Param("id"), userIDStr, userName, input)
	if err != nil {
		switch err.Error() {
		case "invalid destination ID":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "destination not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "you have already reviewed this destination":
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"review": review})
}

func (h *ReviewHandler) GetAllReviews(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reviews, err := h.reviewService.GetAllReviews(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}

func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := h.reviewService.DeleteReview(ctx, c.Param("id")); err != nil {
		switch err.Error() {
		case "invalid review ID":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "review not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}
