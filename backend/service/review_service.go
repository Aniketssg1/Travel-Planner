package service

import (
	"context"
	"errors"
	"math"

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
		return nil, errors.New("invalid destination ID")
	}
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	_, err = s.destRepo.FindByID(ctx, destID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("destination not found")
		}
		return nil, errors.New("database error")
	}

	_, err = s.reviewRepo.FindByUserAndDestination(ctx, userID, destID)
	if err == nil {
		return nil, errors.New("you have already reviewed this destination")
	}

	user, _ := s.userRepo.FindByEmail(ctx, "")
	_ = user

	review := models.Review{
		ID:            primitive.NewObjectID(),
		DestinationID: destID,
		UserID:        userID,
		Rating:        input.Rating,
		Text:          input.Text,
	}

	if err := s.reviewRepo.Create(ctx, review); err != nil {
		return nil, errors.New("failed to submit review")
	}

	s.recalcRating(ctx, destID)
	return &review, nil
}

func (s *ReviewService) SubmitReviewWithUserName(ctx context.Context, destIDStr, userIDStr, userName string, input models.ReviewInput) (*models.Review, error) {
	destID, err := primitive.ObjectIDFromHex(destIDStr)
	if err != nil {
		return nil, errors.New("invalid destination ID")
	}
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	_, err = s.destRepo.FindByID(ctx, destID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("destination not found")
		}
		return nil, errors.New("database error")
	}

	_, err = s.reviewRepo.FindByUserAndDestination(ctx, userID, destID)
	if err == nil {
		return nil, errors.New("you have already reviewed this destination")
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
		return nil, errors.New("failed to submit review")
	}

	s.recalcRating(ctx, destID)
	return &review, nil
}

func (s *ReviewService) GetAllReviews(ctx context.Context) ([]models.Review, error) {
	reviews, err := s.reviewRepo.FindAll(ctx)
	if err != nil {
		return nil, errors.New("database error")
	}
	return reviews, nil
}

func (s *ReviewService) DeleteReview(ctx context.Context, reviewIDStr string) error {
	reviewID, err := primitive.ObjectIDFromHex(reviewIDStr)
	if err != nil {
		return errors.New("invalid review ID")
	}

	review, err := s.reviewRepo.FindByID(ctx, reviewID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("review not found")
		}
		return errors.New("database error")
	}

	if err := s.reviewRepo.Delete(ctx, reviewID); err != nil {
		return errors.New("failed to delete review")
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
