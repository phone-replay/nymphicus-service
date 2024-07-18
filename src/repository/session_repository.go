package repository

import (
	"context"
	"nymphicus-service/src/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type SessionRepository interface {
	SaveActionsToMongo(actions models.Session) error
	UpdateSessionStatusToError(key string) error
}

type sessionRepository struct {
	database *mongo.Database
}

func NewSessionRepository(database *mongo.Database) SessionRepository {
	return &sessionRepository{database: database}
}

func (c *sessionRepository) SaveActionsToMongo(actions models.Session) error {
	collection := c.database.Collection("sessions")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, actions)
	return err
}

func (c *sessionRepository) UpdateSessionStatusToError(key string) error {
	collection := c.database.Collection("sessions")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"key": key}
	update := bson.M{
		"$set": bson.M{
			"status": "Error",
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}
