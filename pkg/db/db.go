package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


var Client *mongo.Client

func Connect(uri string) error {

	// set client options
	clientOptions := options.Client().ApplyURI(uri)

	// connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return err
	}

	Client = client

	// print the Client object
	log.Println("MongoDB Client:", fmt.Sprintf("%+v", Client))

	

	println("Connected to MongoDB!")
	return nil
}

func Disconnect() error {
	err := Client.Disconnect(context.Background())
	if err != nil {
		return err
	}
	println("Disconnected from MongoDB!")
	return nil
}
