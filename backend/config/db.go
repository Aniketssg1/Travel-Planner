package config

import (
	"context"
	"os"
	"strconv"
	"time"

	"backend/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	DB     *mongo.Database
	client *mongo.Client
)

func envInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

// ConnectDB establishes a MongoDB connection with production-grade pooling.
func ConnectDB() {
	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DB_NAME")

	maxPool := uint64(envInt("MONGO_MAX_POOL_SIZE", 200))
	minPool := uint64(envInt("MONGO_MIN_POOL_SIZE", 20))

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetMaxPoolSize(maxPool).
		SetMinPoolSize(minPool).
		SetMaxConnIdleTime(60 * time.Second).
		SetConnectTimeout(10 * time.Second).
		SetServerSelectionTimeout(5 * time.Second).
		SetRetryWrites(true).
		SetRetryReads(true).
		SetCompressors([]string{"snappy", "zstd"})

	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to connect to MongoDB")
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to ping MongoDB")
	}

	DB = client.Database(dbName)
	logger.Log.Info().
		Str("database", dbName).
		Uint64("maxPoolSize", maxPool).
		Uint64("minPoolSize", minPool).
		Msg("Connected to MongoDB")
}

// PingDB checks database connectivity (used by readiness probe).
func PingDB(ctx context.Context) error {
	return client.Ping(ctx, readpref.Primary())
}

// DisconnectDB gracefully closes the MongoDB connection.
func DisconnectDB() {
	if client == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Disconnect(ctx); err != nil {
		logger.Log.Error().Err(err).Msg("Error disconnecting from MongoDB")
	} else {
		logger.Log.Info().Msg("MongoDB connection closed")
	}
}

func GetCollection(name string) *mongo.Collection {
	return DB.Collection("TP_" + name)
}
