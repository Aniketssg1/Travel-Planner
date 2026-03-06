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

type UserRepository struct {
	col *mongo.Collection
}

func NewUserRepository() *UserRepository {
	return &UserRepository{col: config.GetCollection("users")}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := r.col.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	return user, err
}

func (r *UserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (models.User, error) {
	var user models.User
	err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	return user, err
}

func (r *UserRepository) Create(ctx context.Context, user models.User) error {
	if user.ID.IsZero() {
		user.ID = primitive.NewObjectID()
	}
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}
	_, err := r.col.InsertOne(ctx, user)
	return err
}

func (r *UserRepository) CreateEmailIndex(ctx context.Context) {
	r.col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	})
}
