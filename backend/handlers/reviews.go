package handlers

import (
	"net/http"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": "VALIDATION_ERROR"})
		return
	}

	ctx := c.Request.Context()
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
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"review": review})
}

func (h *ReviewHandler) GetAllReviews(c *gin.Context) {
	reviews, err := h.reviewService.GetAllReviews(c.Request.Context())
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}

func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	if err := h.reviewService.DeleteReview(c.Request.Context(), c.Param("id")); err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}
