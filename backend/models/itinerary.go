package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Itinerary struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	UserID       primitive.ObjectID   `bson:"userId" json:"userId"`
	Name         string               `bson:"name" json:"name"`
	StartDate    time.Time            `bson:"startDate" json:"startDate"`
	EndDate      time.Time            `bson:"endDate" json:"endDate"`
	Destinations []primitive.ObjectID `bson:"destinations" json:"destinations"`
	CreatedAt    time.Time            `bson:"createdAt" json:"createdAt"`
}

type ItineraryInput struct {
	Name         string   `json:"name" binding:"required"`
	StartDate    string   `json:"startDate" binding:"required"`
	EndDate      string   `json:"endDate" binding:"required"`
	Destinations []string `json:"destinations"`
}

type ItineraryWithDestinations struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID       primitive.ObjectID `bson:"userId" json:"userId"`
	Name         string             `bson:"name" json:"name"`
	StartDate    time.Time          `bson:"startDate" json:"startDate"`
	EndDate      time.Time          `bson:"endDate" json:"endDate"`
	Destinations []Destination      `bson:"destinationDetails" json:"destinations"`
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`
}
