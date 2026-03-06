package service

import (
	"context"
	"errors"
	"time"

	"backend/models"
	"backend/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ItineraryService struct {
	itinRepo *repository.ItineraryRepository
	destRepo *repository.DestinationRepository
}

func NewItineraryService(itinRepo *repository.ItineraryRepository, destRepo *repository.DestinationRepository) *ItineraryService {
	return &ItineraryService{itinRepo: itinRepo, destRepo: destRepo}
}

type ItineraryResponse struct {
	models.Itinerary
	DestinationDetails []models.Destination `json:"destinationDetails"`
}

func (s *ItineraryService) GetItineraries(ctx context.Context, userIDStr string) ([]ItineraryResponse, error) {
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	itineraries, err := s.itinRepo.FindByUser(ctx, userID)
	if err != nil {
		return nil, errors.New("database error")
	}

	var enriched []ItineraryResponse
	for _, itin := range itineraries {
		dests, err := s.destRepo.FindByIDs(ctx, itin.Destinations)
		if err != nil {
			dests = []models.Destination{}
		}

		// Sort dests to match the order of itin.Destinations (MongoDB $in does NOT preserve order)
		idOrder := make(map[primitive.ObjectID]int)
		for i, oid := range itin.Destinations {
			idOrder[oid] = i
		}
		sortedDests := make([]models.Destination, len(dests))
		for _, d := range dests {
			if idx, ok := idOrder[d.ID]; ok && idx < len(sortedDests) {
				sortedDests[idx] = d
			}
		}
		var finalDests []models.Destination
		for _, d := range sortedDests {
			if !d.ID.IsZero() {
				finalDests = append(finalDests, d)
			}
		}
		if finalDests == nil {
			finalDests = []models.Destination{}
		}

		enriched = append(enriched, ItineraryResponse{itin, finalDests})
	}
	if enriched == nil {
		enriched = []ItineraryResponse{}
	}

	return enriched, nil
}

func (s *ItineraryService) CreateItinerary(ctx context.Context, userIDStr string, input models.ItineraryInput) (*models.Itinerary, error) {
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		return nil, errors.New("invalid startDate format, use YYYY-MM-DD")
	}
	endDate, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		return nil, errors.New("invalid endDate format, use YYYY-MM-DD")
	}
	if endDate.Before(startDate) {
		return nil, errors.New("endDate must be after startDate")
	}

	var destIDs []primitive.ObjectID
	for _, idStr := range input.Destinations {
		oid, err := primitive.ObjectIDFromHex(idStr)
		if err == nil {
			destIDs = append(destIDs, oid)
		}
	}
	if destIDs == nil {
		destIDs = []primitive.ObjectID{}
	}

	itin := models.Itinerary{
		UserID:       userID,
		Name:         input.Name,
		StartDate:    startDate,
		EndDate:      endDate,
		Destinations: destIDs,
	}

	created, err := s.itinRepo.Create(ctx, itin)
	if err != nil {
		return nil, errors.New("failed to create itinerary")
	}
	return &created, nil
}

func (s *ItineraryService) UpdateItinerary(ctx context.Context, itinIDStr, userIDStr string, input models.ItineraryInput) error {
	itinID, err := primitive.ObjectIDFromHex(itinIDStr)
	if err != nil {
		return errors.New("invalid itinerary ID")
	}
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return errors.New("invalid user ID")
	}

	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		return errors.New("invalid startDate format, use YYYY-MM-DD")
	}
	endDate, err := time.Parse("2006-01-02", input.EndDate)
	if err != nil {
		return errors.New("invalid endDate format, use YYYY-MM-DD")
	}
	if endDate.Before(startDate) {
		return errors.New("endDate must be after startDate")
	}

	var destIDs []primitive.ObjectID
	for _, idStr := range input.Destinations {
		oid, err := primitive.ObjectIDFromHex(idStr)
		if err == nil {
			destIDs = append(destIDs, oid)
		}
	}
	if destIDs == nil {
		destIDs = []primitive.ObjectID{}
	}

	update := bson.M{
		"name":         input.Name,
		"startDate":    startDate,
		"endDate":      endDate,
		"destinations": destIDs,
	}

	if err := s.itinRepo.Update(ctx, itinID, userID, update); err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("itinerary not found or access denied")
		}
		return errors.New("failed to update itinerary")
	}
	return nil
}

func (s *ItineraryService) DeleteItinerary(ctx context.Context, itinIDStr, userIDStr string) error {
	itinID, err := primitive.ObjectIDFromHex(itinIDStr)
	if err != nil {
		return errors.New("invalid itinerary ID")
	}
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		return errors.New("invalid user ID")
	}

	if err := s.itinRepo.Delete(ctx, itinID, userID); err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("itinerary not found or access denied")
		}
		return errors.New("failed to delete itinerary")
	}
	return nil
}
