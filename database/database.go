package database

import (
	"context"
	"nymphicus-service/config"
	"nymphicus-service/pkg/logger"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoConnectTimeout = 10 * time.Second
	mongoPingTimeout    = 2 * time.Second
)

func ConnectionDatabase(c *config.Config, logger logger.Logger) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), mongoConnectTimeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI(c.MongoDB.MongoURI)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Fatalf("Failed to connect to MongoDB: %v", err)
		return nil, err
	}

	pingCtx, pingCancel := context.WithTimeout(context.Background(), mongoPingTimeout)
	defer pingCancel()

	if err := client.Ping(pingCtx, nil); err != nil {
		logger.Fatalf("Failed to ping MongoDB: %v", err)
		return nil, err
	}

	logger.Infof("Connected to MongoDB!")

	if c.Server.Mode == "production" {
		logger.Infof("Running in production mode.")
	} else {
		logger.Infof("Running in development mode.")
	}

	database := client.Database(c.MongoDB.Database)

	return database, nil
}
