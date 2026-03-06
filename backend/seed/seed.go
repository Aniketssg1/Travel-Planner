package seed

import (
	"context"
	"log"
	"time"

	"backend/config"
	"backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func SeedData() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	seedAdmin(ctx)
	seedDestinations(ctx)
	createIndexes(ctx)
}

func seedAdmin(ctx context.Context) {
	col := config.GetCollection("users")

	col.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	})

	var existing models.User
	err := col.FindOne(ctx, bson.M{"email": "admin@travel.com"}).Decode(&existing)
	if err == nil {
		log.Println("Admin user already exists, skipping seed")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Admin@123"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash admin password: %v", err)
		return
	}

	admin := models.User{
		ID:        primitive.NewObjectID(),
		Name:      "Admin",
		Email:     "admin@travel.com",
		Password:  string(hashedPassword),
		Role:      "admin",
		CreatedAt: time.Now(),
	}

	_, err = col.InsertOne(ctx, admin)
	if err != nil {
		log.Printf("Failed to seed admin: %v", err)
		return
	}
	log.Println("Admin user seeded: admin@travel.com / Admin@123")
}

func seedDestinations(ctx context.Context) {
	col := config.GetCollection("destinations")

	count, _ := col.CountDocuments(ctx, bson.M{})
	if count > 0 {
		log.Println("Destinations already exist, skipping seed")
		return
	}

	destinations := []interface{}{
		models.Destination{
			ID:            primitive.NewObjectID(),
			Name:          "Paris",
			Country:       "France",
			Description:   "The City of Light, known for the Eiffel Tower, world-class cuisine, and romantic atmosphere.",
			ImageURL:      "https://images.unsplash.com/photo-1502602898657-3e91760cbb34?w=800",
			AverageRating: 0,
			ReviewCount:   0,
			CreatedAt:     time.Now(),
		},
		models.Destination{
			ID:            primitive.NewObjectID(),
			Name:          "Tokyo",
			Country:       "Japan",
			Description:   "A vibrant metropolis blending ultramodern and traditional, from neon-lit skyscrapers to historic temples.",
			ImageURL:      "https://images.unsplash.com/photo-1540959733332-eab4deabeeaf?w=800",
			AverageRating: 0,
			ReviewCount:   0,
			CreatedAt:     time.Now(),
		},
		models.Destination{
			ID:            primitive.NewObjectID(),
			Name:          "New York City",
			Country:       "USA",
			Description:   "The Big Apple — iconic skyline, Times Square, Central Park, and endless cultural experiences.",
			ImageURL:      "https://images.unsplash.com/photo-1496442226666-8d4d0e62e6e9?w=800",
			AverageRating: 0,
			ReviewCount:   0,
			CreatedAt:     time.Now(),
		},
		models.Destination{
			ID:            primitive.NewObjectID(),
			Name:          "Bali",
			Country:       "Indonesia",
			Description:   "Tropical paradise with stunning rice terraces, ancient temples, and beautiful beaches.",
			ImageURL:      "https://images.unsplash.com/photo-1537996194471-e657df975ab4?w=800",
			AverageRating: 0,
			ReviewCount:   0,
			CreatedAt:     time.Now(),
		},
		models.Destination{
			ID:            primitive.NewObjectID(),
			Name:          "Rome",
			Country:       "Italy",
			Description:   "The Eternal City, home to the Colosseum, Vatican, and some of the world's finest food.",
			ImageURL:      "https://images.unsplash.com/photo-1552832230-c0197dd311b5?w=800",
			AverageRating: 0,
			ReviewCount:   0,
			CreatedAt:     time.Now(),
		},
		models.Destination{
			ID:            primitive.NewObjectID(),
			Name:          "Sydney",
			Country:       "Australia",
			Description:   "Stunning harbour city with the iconic Opera House, beautiful beaches, and vibrant culture.",
			ImageURL:      "https://images.unsplash.com/photo-1506973035872-a4ec16b8e8d9?w=800",
			AverageRating: 0,
			ReviewCount:   0,
			CreatedAt:     time.Now(),
		},
	}

	_, err := col.InsertMany(ctx, destinations)
	if err != nil {
		log.Printf("Failed to seed destinations: %v", err)
		return
	}
	log.Println("Sample destinations seeded successfully")
}

func createIndexes(ctx context.Context) {
	reviewCol := config.GetCollection("reviews")
	reviewCol.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "destinationId", Value: 1}, {Key: "userId", Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	itinCol := config.GetCollection("itineraries")
	itinCol.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{"userId": 1},
	})

	destCol := config.GetCollection("destinations")
	destCol.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "country", Value: 1}, {Key: "name", Value: 1}},
	})

	log.Println("Database indexes created")
}
