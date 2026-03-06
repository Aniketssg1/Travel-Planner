package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Destination struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `bson:"name" json:"name"`
	Country       string             `bson:"country" json:"country"`
	Description   string             `bson:"description" json:"description"`
	ImageURL      string             `bson:"imageUrl" json:"imageUrl"`
	AverageRating float64            `bson:"averageRating" json:"averageRating"`
	ReviewCount   int                `bson:"reviewCount" json:"reviewCount"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
}

type DestinationInput struct {
	Name        string `json:"name" binding:"required"`
	Country     string `json:"country" binding:"required"`
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
}
