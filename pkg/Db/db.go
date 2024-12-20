package db

import (
	"context"
	"fmt"
	"log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User represents a user model
type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Firstname        string            `json:"firstname" bson:"firstname"`
	Lastname         string            `json:"lastname" bson:"lastname"`
	Email            string            `json:"email" bson:"email"`
	Phone            string            `json:"phone" bson:"phone"`
	Password         string            `json:"-" bson:"password"` // Hide password from JSON responses
	CreatedAt        primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt        primitive.DateTime `json:"updated_at" bson:"updated_at"`
	Orders           []string          `json:"orders" bson:"orders"`
	Cart             []string          `json:"cart" bson:"cart"`
	IsActive         bool              `json:"is_active" bson:"is_active"`
	Roles            []string          `json:"roles" bson:"roles"`
	ProfilePictureURL string           `json:"profile_picture_url" bson:"profile_picture_url,omitempty"`
	LastLogin        primitive.DateTime `json:"last_login,omitempty" bson:"last_login,omitempty"`
} 

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
