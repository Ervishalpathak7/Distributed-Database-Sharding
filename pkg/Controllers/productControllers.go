package Controllers

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/Ervishalpathak7/Distributed-Database-Sharding/pkg/Db"
	"github.com/Ervishalpathak7/Distributed-Database-Sharding/pkg/Schemas"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetProduct retrieves a product by its productId
func GetProduct(c *gin.Context) {
	productId := c.Param("productId")
	objectId, err := primitive.ObjectIDFromHex(productId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid productId"})
		return
	}

	collection := db.Client.Database("shardDB").Collection("products")

	var product Schema.Product

	err = collection.FindOne(c.Request.Context(), bson.M{"_id": objectId}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": product})
}


// ListProducts retrieves a paginated list of products with optional filters
func ListProducts(c *gin.Context) {
	var page int64 = 1
	var limit int64 = 10
	var err error

	// Parse pagination parameters
	if pageStr := c.Query("page"); pageStr != "" {
		page, err = strconv.ParseInt(pageStr, 10, 64)
		if err != nil || page < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
			return
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err = strconv.ParseInt(limitStr, 10, 64)
		if err != nil || limit < 1 || limit > 100 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit. Must be between 1 and 100"})
			return
		}
	}

	// Build filter
	filter := bson.M{}

	if category := c.Query("category"); category != "" {
		filter["category"] = category
	}

	if minPrice := c.Query("minPrice"); minPrice != "" {
		price, err := strconv.ParseFloat(minPrice, 64)
		if err == nil {
			filter["price"] = bson.M{"$gte": price}
		}
	}

	if maxPrice := c.Query("maxPrice"); maxPrice != "" {
		price, err := strconv.ParseFloat(maxPrice, 64)
		if err == nil {
			if _, ok := filter["price"]; ok {
				filter["price"].(bson.M)["$lte"] = price
			} else {
				filter["price"] = bson.M{"$lte": price}
			}
		}
	}

	if inStock := c.Query("inStock"); inStock == "true" {
		filter["stock_quantity"] = bson.M{"$gt": 0}
	}

	collection := db.Client.Database("shardDB").Collection("products")
	
	// Set up options for pagination and sorting
	skip := (page - 1) * limit
	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.M{"created_at": -1})

	cursor, err := collection.Find(c.Request.Context(), filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(c.Request.Context())

	var products []Schema.Product
	if err := cursor.All(c.Request.Context(), &products); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get total count for pagination
	total, err := collection.CountDocuments(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"products": products,
			"total":    total,
			"page":     page,
			"limit":    limit,
			"pages":    math.Ceil(float64(total) / float64(limit)),
		},
	})
}

// CreateProduct creates a new product in MongoDB
func CreateProduct(c *gin.Context) {
	var product Schema.Product
	product.ID = primitive.NewObjectID()
	product.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	product.UpdatedAt = product.CreatedAt

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate required fields
	if product.Name == "" || product.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and price are required. Price must be greater than 0"})
		return
	}

	collection := db.Client.Database("shardDB").Collection("products")

	// Check for existing product with same name
	var existingProduct Schema.Product
	err := collection.FindOne(c.Request.Context(), bson.M{"name": product.Name}).Decode(&existingProduct)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Product with this name already exists"})
		return
	} else if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result, err := collection.InsertOne(c.Request.Context(), product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"id":      result.InsertedID,
	})
}

// UpdateProduct updates a product by its productId
func UpdateProduct(c *gin.Context) {
	productId := c.Param("productId")

	var updates struct {
		Name          string   `json:"name"`
		Description   string   `json:"description"`
		Price         float64  `json:"price"`
		Category      string   `json:"category"`
		Images        []string `json:"images"`
		StockQuantity int      `json:"stock_quantity"`
		IsActive      *bool    `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectId, err := primitive.ObjectIDFromHex(productId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid productId"})
		return
	}

	collection := db.Client.Database("shardDB").Collection("products")

	// Check name uniqueness if name is being updated
	if updates.Name != "" {
		var existingProduct Schema.Product
		err := collection.FindOne(c.Request.Context(), bson.M{
			"_id":  bson.M{"$ne": objectId},
			"name": updates.Name,
		}).Decode(&existingProduct)
		
		if err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Product with this name already exists"})
			return
		} else if err != mongo.ErrNoDocuments {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	updateFields := bson.M{
		"updated_at": primitive.NewDateTimeFromTime(time.Now()),
	}

	if updates.Name != "" {
		updateFields["name"] = updates.Name
	}
	if updates.Description != "" {
		updateFields["description"] = updates.Description
	}
	if updates.Price > 0 {
		updateFields["price"] = updates.Price
	}
	if updates.Category != "" {
		updateFields["category"] = updates.Category
	}
	if updates.Images != nil {
		updateFields["images"] = updates.Images
	}
	if updates.StockQuantity >= 0 {
		updateFields["stock_quantity"] = updates.StockQuantity
	}
	if updates.IsActive != nil {
		updateFields["is_active"] = *updates.IsActive
	}

	result, err := collection.UpdateOne(
		c.Request.Context(),
		bson.M{"_id": objectId},
		bson.M{"$set": updateFields},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found or no changes made"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Product updated successfully"})
}

// DeleteProduct deletes a product by its productId
func DeleteProduct(c *gin.Context) {
	productId := c.Param("productId")

	objectId, err := primitive.ObjectIDFromHex(productId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid productId"})
		return
	}

	collection := db.Client.Database("shardDB").Collection("products")
	result, err := collection.DeleteOne(c.Request.Context(), bson.M{"_id": objectId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Product deleted successfully"})
}