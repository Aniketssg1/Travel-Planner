package repository

import (
	"context"
	"time"

	"backend/config"
	"backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ItineraryRepository struct {
	col *mongo.Collection
}

func NewItineraryRepository() *ItineraryRepository {
	return &ItineraryRepository{col: config.GetCollection("itineraries")}
}

func (r *ItineraryRepository) FindByUser(ctx context.Context, userID primitive.ObjectID) ([]models.Itinerary, error) {
	cursor, err := r.col.Find(ctx, bson.M{"userId": userID}, options.Find().SetSort(bson.M{"createdAt": -1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var itineraries []models.Itinerary
	if err := cursor.All(ctx, &itineraries); err != nil {
		return nil, err
	}
	if itineraries == nil {
		itineraries = []models.Itinerary{}
	}
	return itineraries, nil
}

func (r *ItineraryRepository) Create(ctx context.Context, itin models.Itinerary) (models.Itinerary, error) {
	if itin.ID.IsZero() {
		itin.ID = primitive.NewObjectID()
	}
	if itin.CreatedAt.IsZero() {
		itin.CreatedAt = time.Now()
	}
	_, err := r.col.InsertOne(ctx, itin)
	return itin, err
}

func (r *ItineraryRepository) Update(ctx context.Context, id, userID primitive.ObjectID, update bson.M) error {
	result, err := r.col.UpdateOne(ctx, bson.M{"_id": id, "userId": userID}, bson.M{"$set": update})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *ItineraryRepository) Delete(ctx context.Context, id, userID primitive.ObjectID) error {
	result, err := r.col.DeleteOne(ctx, bson.M{"_id": id, "userId": userID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *ItineraryRepository) CreateIndexes(ctx context.Context) {
	r.col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{"userId": 1},
	})
}
