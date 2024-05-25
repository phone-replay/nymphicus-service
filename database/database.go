package database

import (
	"context"
	"fmt"
	"log"
	"nymphicus-service/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoConnectTimeout = 10 * time.Second
	mongoPingTimeout    = 2 * time.Second
)

func ConnectionDatabase(c *config.Config) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), mongoConnectTimeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI(c.MongoDB.MongoURI)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	pingCtx, pingCancel := context.WithTimeout(context.Background(), mongoPingTimeout)
	defer pingCancel()

	if err := client.Ping(pingCtx, nil); err != nil {
		log.Fatal(err)
		return nil, err
	}

	fmt.Println("Conectado ao MongoDB!")

	if c.Server.Mode == "production" {
		fmt.Println("Está em ambiente de produção.")
	} else {
		fmt.Println("Está em ambiente de desenvolvimento.")
	}

	return client, nil
}
