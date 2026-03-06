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

type ReviewRepository struct {
	col *mongo.Collection
}

func NewReviewRepository() *ReviewRepository {
	return &ReviewRepository{col: config.GetCollection("reviews")}
}

func (r *ReviewRepository) FindByDestination(ctx context.Context, destID primitive.ObjectID) ([]models.Review, error) {
	cursor, err := r.col.Find(ctx, bson.M{"destinationId": destID}, options.Find().SetSort(bson.M{"createdAt": -1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reviews []models.Review
	if err := cursor.All(ctx, &reviews); err != nil {
		return nil, err
	}
	if reviews == nil {
		reviews = []models.Review{}
	}
	return reviews, nil
}

func (r *ReviewRepository) FindAll(ctx context.Context) ([]models.Review, error) {
	cursor, err := r.col.Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"createdAt": -1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reviews []models.Review
	if err := cursor.All(ctx, &reviews); err != nil {
		return nil, err
	}
	if reviews == nil {
		reviews = []models.Review{}
	}
	return reviews, nil
}

func (r *ReviewRepository) FindByUserAndDestination(ctx context.Context, userID, destID primitive.ObjectID) (models.Review, error) {
	var review models.Review
	err := r.col.FindOne(ctx, bson.M{"destinationId": destID, "userId": userID}).Decode(&review)
	return review, err
}

func (r *ReviewRepository) FindByID(ctx context.Context, id primitive.ObjectID) (models.Review, error) {
	var review models.Review
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&review)
	return review, err
}

func (r *ReviewRepository) Create(ctx context.Context, review models.Review) error {
	if review.ID.IsZero() {
		review.ID = primitive.NewObjectID()
	}
	if review.CreatedAt.IsZero() {
		review.CreatedAt = time.Now()
	}
	_, err := r.col.InsertOne(ctx, review)
	return err
}

func (r *ReviewRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *ReviewRepository) DeleteByDestination(ctx context.Context, destID primitive.ObjectID) error {
	_, err := r.col.DeleteMany(ctx, bson.M{"destinationId": destID})
	return err
}

func (r *ReviewRepository) AggregateRating(ctx context.Context, destID primitive.ObjectID) (float64, int, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"destinationId": destID}}},
		{{Key: "$group", Value: bson.M{
			"_id":        nil,
			"avgRating":  bson.M{"$avg": "$rating"},
			"totalCount": bson.M{"$sum": 1},
		}}},
	}

	cursor, err := r.col.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, 0, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return 0, 0, err
	}

	if len(results) == 0 {
		return 0, 0, nil
	}

	avgRating := results[0]["avgRating"].(float64)
	totalCount := results[0]["totalCount"].(int32)
	return avgRating, int(totalCount), nil
}

func (r *ReviewRepository) CreateIndexes(ctx context.Context) {
	r.col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "destinationId", Value: 1}, {Key: "userId", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
}
