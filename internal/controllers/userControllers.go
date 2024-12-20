package userControllers

import (
	"encoding/json"
	"net/http"
	"github.com/Ervishalpathak7/Distributed-Database-Sharding/pkg/db"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// User represents a user model
type User struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `json:"name" bson:"name"`
	Email string             `json:"email" bson:"email"`
}

// MarshalJSON customizes the way the User struct is returned as JSON
func (u *User) MarshalJSON() ([]byte, error) {
	type Alias User
	return json.Marshal(&struct {
		ID string `json:"id"` // Custom field for ID as a string
		*Alias
	}{
		ID:    u.ID.Hex(), // Convert the ObjectID to string (hex format)
		Alias: (*Alias)(u),
	})
}

// GetUsers retrieves a user by their userId
func GetUsers(c *gin.Context) {
	userId := c.Param("userId")


	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userId"})
		return
	}

	collection := db.Client.Database("shardDB").Collection("users")
	var user User

	err = collection.FindOne(c.Request.Context(), bson.M{"_id": objectId}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": user})
}

// CreateUser creates a new user in MongoDB
func CreateUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Name == "" || user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and Email are required"})
		return
	}

	collection := db.Client.Database("shardDB").Collection("users")
	result, err := collection.InsertOne(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.ID = result.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": user})
}

// UpdateUser updates a user by their userId
func UpdateUser(c *gin.Context) {
	userId := c.Param("userId")
	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Name == "" || user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and Email are required"})
		return
	}

	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userId"})
		return
	}

	collection := db.Client.Database("shardDB").Collection("users")
	update := bson.M{
		"$set": bson.M{
			"name":  user.Name,
			"email": user.Email,
		},
	}

	result, err := collection.UpdateOne(c.Request.Context(), bson.M{"_id": objectId}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found or no changes made"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User updated successfully"})
}

// DeleteUser deletes a user by their userId
func DeleteUser(c *gin.Context) {
	userId := c.Param("userId")


	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userId"})
		return
	}

	collection := db.Client.Database("shardDB").Collection("users")
	result, err := collection.DeleteOne(c.Request.Context(), bson.M{"_id": objectId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User deleted successfully"})
}
