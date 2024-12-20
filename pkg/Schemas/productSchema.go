package Schema

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string            `json:"name" bson:"name"`
	Description   string            `json:"description" bson:"description"`
	Price         float64           `json:"price" bson:"price"`
	Category      string            `json:"category" bson:"category"`
	Images        []string          `json:"images" bson:"images"`
	StockQuantity int               `json:"stock_quantity" bson:"stock_quantity"`
	IsActive      bool              `json:"is_active" bson:"is_active"`
	CreatedAt     primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt     primitive.DateTime `json:"updated_at" bson:"updated_at"`
}

func NewProduct() *Product {
	currentTime := primitive.NewDateTimeFromTime(time.Now())
	return &Product{
		ID:            primitive.NewObjectID(),
		IsActive:      true,
		CreatedAt:     currentTime,
		UpdatedAt:     currentTime,
		Images:        []string{},
		StockQuantity: 0,
	}
}