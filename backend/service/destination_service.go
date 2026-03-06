package service

import (
	"context"
	"errors"
	"strconv"

	"backend/models"
	"backend/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DestinationService struct {
	destRepo   *repository.DestinationRepository
	reviewRepo *repository.ReviewRepository
}

func NewDestinationService(destRepo *repository.DestinationRepository, reviewRepo *repository.ReviewRepository) *DestinationService {
	return &DestinationService{destRepo: destRepo, reviewRepo: reviewRepo}
}

type ListFilters struct {
	Name      string
	Country   string
	MinRating string
	Page      string
	Limit     string
}

type ListResult struct {
	Destinations []models.Destination
	Pagination   repository.PaginationResult
}

func (s *DestinationService) ListDestinations(ctx context.Context, filters ListFilters) (*ListResult, error) {
	page, _ := strconv.Atoi(filters.Page)
	limit, _ := strconv.Atoi(filters.Limit)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}

	filter := bson.M{}
	if filters.Name != "" {
		filter["name"] = bson.M{"$regex": filters.Name, "$options": "i"}
	}
	if filters.Country != "" {
		filter["country"] = bson.M{"$regex": filters.Country, "$options": "i"}
	}
	if filters.MinRating != "" {
		rating, err := strconv.ParseFloat(filters.MinRating, 64)
		if err == nil {
			filter["averageRating"] = bson.M{"$gte": rating}
		}
	}

	destinations, pagination, err := s.destRepo.FindAll(ctx, filter, page, limit)
	if err != nil {
		return nil, errors.New("database error")
	}

	return &ListResult{Destinations: destinations, Pagination: pagination}, nil
}

func (s *DestinationService) GetDestination(ctx context.Context, idStr string) (*models.Destination, []models.Review, error) {
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return nil, nil, errors.New("invalid destination ID")
	}

	dest, err := s.destRepo.FindByID(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil, errors.New("destination not found")
		}
		return nil, nil, errors.New("database error")
	}

	reviews, err := s.reviewRepo.FindByDestination(ctx, id)
	if err != nil {
		return nil, nil, errors.New("failed to fetch reviews")
	}

	return &dest, reviews, nil
}

func (s *DestinationService) CreateDestination(ctx context.Context, input models.DestinationInput) (*models.Destination, error) {
	dest := models.Destination{
		Name:          input.Name,
		Country:       input.Country,
		Description:   input.Description,
		ImageURL:      input.ImageURL,
		AverageRating: 0,
		ReviewCount:   0,
	}

	created, err := s.destRepo.Create(ctx, dest)
	if err != nil {
		return nil, errors.New("failed to create destination")
	}
	return &created, nil
}

func (s *DestinationService) UpdateDestination(ctx context.Context, idStr string, input models.DestinationInput) error {
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return errors.New("invalid destination ID")
	}

	update := bson.M{
		"name":        input.Name,
		"country":     input.Country,
		"description": input.Description,
		"imageUrl":    input.ImageURL,
	}

	if err := s.destRepo.Update(ctx, id, update); err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("destination not found")
		}
		return errors.New("failed to update destination")
	}
	return nil
}

func (s *DestinationService) DeleteDestination(ctx context.Context, idStr string) error {
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return errors.New("invalid destination ID")
	}

	if err := s.destRepo.Delete(ctx, id); err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("destination not found")
		}
		return errors.New("failed to delete destination")
	}

	s.reviewRepo.DeleteByDestination(ctx, id)
	return nil
}
