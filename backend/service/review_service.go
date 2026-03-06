package service

import (
	"context"
	"math"

	apperrors "backend/errors"
	"backend/models"
	"backend/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReviewService struct {
	reviewRepo *repository.ReviewRepository
	destRepo   *repository.DestinationRepository
	userRepo   *repository.UserRepository
}

func NewReviewService(
	reviewRepo *repository.ReviewRepository,
	destRepo *repository.DestinationRepository,
	userRepo *repository.UserRepository,
) *ReviewService {
	return &ReviewService{reviewRepo: reviewRepo, destRepo: destRepo, userRepo: userRepo}
}

func (s *ReviewService) SubmitReview(ctx context.Context, destIDStr, userIDStr string, input models.ReviewInput) (*models.Review, error) {
	destID, err := primitive.ObjectIDFromHex(destIDStr)
	if err != nil {
		return nil, apperrors.BadRequest("invalid destination ID")
	}
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return nil, apperrors.BadRequest("invalid user ID")
	}

	_, err = s.destRepo.FindByID(ctx, destID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, apperrors.NotFound("destination not found")
		}
		return nil, apperrors.Internal("database error")
	}

	_, err = s.reviewRepo.FindByUserAndDestination(ctx, userID, destID)
	if err == nil {
		return nil, apperrors.Conflict("you have already reviewed this destination")
	}

	review := models.Review{
		ID:            primitive.NewObjectID(),
		DestinationID: destID,
		UserID:        userID,
		Rating:        input.Rating,
		Text:          input.Text,
	}

	if err := s.reviewRepo.Create(ctx, review); err != nil {
		return nil, apperrors.Internal("failed to submit review")
	}

	s.recalcRating(ctx, destID)
	return &review, nil
}

func (s *ReviewService) SubmitReviewWithUserName(ctx context.Context, destIDStr, userIDStr, userName string, input models.ReviewInput) (*models.Review, error) {
	destID, err := primitive.ObjectIDFromHex(destIDStr)
	if err != nil {
		return nil, apperrors.BadRequest("invalid destination ID")
	}
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return nil, apperrors.BadRequest("invalid user ID")
	}

	_, err = s.destRepo.FindByID(ctx, destID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, apperrors.NotFound("destination not found")
		}
		return nil, apperrors.Internal("database error")
	}

	_, err = s.reviewRepo.FindByUserAndDestination(ctx, userID, destID)
	if err == nil {
		return nil, apperrors.Conflict("you have already reviewed this destination")
	}

	review := models.Review{
		ID:            primitive.NewObjectID(),
		DestinationID: destID,
		UserID:        userID,
		UserName:      userName,
		Rating:        input.Rating,
		Text:          input.Text,
	}

	if err := s.reviewRepo.Create(ctx, review); err != nil {
		return nil, apperrors.Internal("failed to submit review")
	}

	s.recalcRating(ctx, destID)
	return &review, nil
}

func (s *ReviewService) GetAllReviews(ctx context.Context) ([]models.Review, error) {
	reviews, err := s.reviewRepo.FindAll(ctx)
	if err != nil {
		return nil, apperrors.Internal("database error")
	}
	return reviews, nil
}

func (s *ReviewService) DeleteReview(ctx context.Context, reviewIDStr string) error {
	reviewID, err := primitive.ObjectIDFromHex(reviewIDStr)
	if err != nil {
		return apperrors.BadRequest("invalid review ID")
	}

	review, err := s.reviewRepo.FindByID(ctx, reviewID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return apperrors.NotFound("review not found")
		}
		return apperrors.Internal("database error")
	}

	if err := s.reviewRepo.Delete(ctx, reviewID); err != nil {
		return apperrors.Internal("failed to delete review")
	}

	s.recalcRating(ctx, review.DestinationID)
	return nil
}

func (s *ReviewService) recalcRating(ctx context.Context, destID primitive.ObjectID) {
	avgRating, count, err := s.reviewRepo.AggregateRating(ctx, destID)
	if err != nil {
		return
	}
	if count == 0 {
		s.destRepo.UpdateRating(ctx, destID, 0, 0)
	} else {
		rounded := math.Round(avgRating*10) / 10
		s.destRepo.UpdateRating(ctx, destID, rounded, count)
	}
}
