package userControllers

import (
	"net/http"
	"time"
	 "github.com/Ervishalpathak7/Distributed-Database-Sharding/pkg/Schemas"
	 "github.com/Ervishalpathak7/Distributed-Database-Sharding/pkg/Utils"
	"github.com/Ervishalpathak7/Distributed-Database-Sharding/pkg/Db"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetUsers retrieves a user by their userId
func GetUsers(c *gin.Context) {
	userId := c.Param("userId")

	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userId"})
		return
	}

	collection := db.Client.Database("shardDB").Collection("users")
	var user userSchema.User

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
	user := userSchema.NewUser()

	if err := c.ShouldBindJSON(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Firstname == "" || user.Email == "" || user.Password == "" || user.Phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Firstname, Email, Password and Phone are required"})
		return
	}

	collection := db.Client.Database("shardDB").Collection("users")

	// Check if the user already exists
	var existingUser userSchema.User
	err := collection.FindOne(c.Request.Context(), bson.M{
		"$or": []bson.M{
			{"email": user.Email},
			{"phone": user.Phone},
		},
	}).Decode(&existingUser)

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this email or phone already exists"})
		return
	} else if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := userUtils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user.Password = hashedPassword
	user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	result, err := collection.InsertOne(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"id":      result.InsertedID,
	})
}

// UpdateUser updates a user by their userId
func UpdateUser(c *gin.Context) {
    userId := c.Param("userId")

    var updates struct {
        Firstname         string `json:"firstname"`
        Lastname          string `json:"lastname"`
        Email             string `json:"email"`
        Phone             string `json:"phone"`
        ProfilePictureURL string `json:"profile_picture_url"`
    }

    if err := c.ShouldBindJSON(&updates); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    objectId, err := primitive.ObjectIDFromHex(userId)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userId"})
        return
    }

    collection := db.Client.Database("shardDB").Collection("users")

    // Check email/phone uniqueness if they're being updated
    if updates.Email != "" || updates.Phone != "" {
        var existingUser userSchema.User
        filter := bson.M{
            "$and": []bson.M{
                {"_id": bson.M{"$ne": objectId}},
                {"$or": []bson.M{
                    {"email": updates.Email},
                    {"phone": updates.Phone},
                }},
            },
        }
        
        err := collection.FindOne(c.Request.Context(), filter).Decode(&existingUser)
        if err == nil {
            c.JSON(http.StatusConflict, gin.H{"error": "Email or phone already in use by another user"})
            return
        } else if err != mongo.ErrNoDocuments {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
    }

    updateFields := bson.M{
        "updated_at": primitive.NewDateTimeFromTime(time.Now()),
    }

    if updates.Firstname != "" {
        updateFields["firstname"] = updates.Firstname
    }
    if updates.Lastname != "" {
        updateFields["lastname"] = updates.Lastname
    }
    if updates.Email != "" {
        updateFields["email"] = updates.Email
    }
    if updates.Phone != "" {
        updateFields["phone"] = updates.Phone
    }
    if updates.ProfilePictureURL != "" {
        updateFields["profile_picture_url"] = updates.ProfilePictureURL
    }

    result, err := collection.UpdateOne(
        c.Request.Context(),
        bson.M{"_id": objectId},
        bson.M{"$set": updateFields},
    )

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "User deleted successfully"})
}
