package repository

import (
	"context"
	"math"
	"time"

	"backend/config"
	"backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DestinationRepository struct {
	col *mongo.Collection
}

func NewDestinationRepository() *DestinationRepository {
	return &DestinationRepository{col: config.GetCollection("destinations")}
}

type PaginationResult struct {
	Total      int64
	Page       int
	Limit      int
	TotalPages int
}

func (r *DestinationRepository) FindAll(ctx context.Context, filter bson.M, page, limit int) ([]models.Destination, PaginationResult, error) {
	total, err := r.col.CountDocuments(ctx, filter)
	if err != nil {
		return nil, PaginationResult{}, err
	}

	skip := int64((page - 1) * limit)
	opts := options.Find().SetSkip(skip).SetLimit(int64(limit)).SetSort(bson.M{"createdAt": -1})
	cursor, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, PaginationResult{}, err
	}
	defer cursor.Close(ctx)

	var destinations []models.Destination
	if err := cursor.All(ctx, &destinations); err != nil {
		return nil, PaginationResult{}, err
	}
	if destinations == nil {
		destinations = []models.Destination{}
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	pagination := PaginationResult{
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	return destinations, pagination, nil
}

func (r *DestinationRepository) FindByID(ctx context.Context, id primitive.ObjectID) (models.Destination, error) {
	var dest models.Destination
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&dest)
	return dest, err
}

func (r *DestinationRepository) FindByIDs(ctx context.Context, ids []primitive.ObjectID) ([]models.Destination, error) {
	if len(ids) == 0 {
		return []models.Destination{}, nil
	}
	cursor, err := r.col.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var dests []models.Destination
	if err := cursor.All(ctx, &dests); err != nil {
		return nil, err
	}
	if dests == nil {
		dests = []models.Destination{}
	}
	return dests, nil
}

func (r *DestinationRepository) Create(ctx context.Context, dest models.Destination) (models.Destination, error) {
	if dest.ID.IsZero() {
		dest.ID = primitive.NewObjectID()
	}
	if dest.CreatedAt.IsZero() {
		dest.CreatedAt = time.Now()
	}
	_, err := r.col.InsertOne(ctx, dest)
	return dest, err
}

func (r *DestinationRepository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	result, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *DestinationRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.col.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *DestinationRepository) UpdateRating(ctx context.Context, id primitive.ObjectID, avgRating float64, reviewCount int) error {
	_, err := r.col.UpdateOne(ctx, bson.M{"_id": id}, bson.M{
		"$set": bson.M{
			"averageRating": avgRating,
			"reviewCount":   reviewCount,
		},
	})
	return err
}

func (r *DestinationRepository) CountAll(ctx context.Context) (int64, error) {
	return r.col.CountDocuments(ctx, bson.M{})
}

func (r *DestinationRepository) InsertMany(ctx context.Context, docs []interface{}) error {
	_, err := r.col.InsertMany(ctx, docs)
	return err
}

func (r *DestinationRepository) CreateIndexes(ctx context.Context) {
	r.col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "country", Value: 1}, {Key: "name", Value: 1}},
	})
}
