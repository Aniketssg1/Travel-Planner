package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Review struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	DestinationID primitive.ObjectID `bson:"destinationId" json:"destinationId"`
	UserID        primitive.ObjectID `bson:"userId" json:"userId"`
	UserName      string             `bson:"userName" json:"userName"`
	Rating        int                `bson:"rating" json:"rating"`
	Text          string             `bson:"text" json:"text"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
}

type ReviewInput struct {
	Rating int    `json:"rating" binding:"required,min=1,max=5"`
	Text   string `json:"text" binding:"required"`
}
