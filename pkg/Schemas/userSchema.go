package Schema


import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Firstname        string            `json:"firstname" bson:"firstname"`
	Lastname         string            `json:"lastname" bson:"lastname"`
	Email            string            `json:"email" bson:"email"`
	Phone            string            `json:"phone" bson:"phone"`
	Password         string            `json:"-" bson:"password"` 
	CreatedAt        primitive.DateTime `json:"created_at" bson:"created_at"`
	UpdatedAt        primitive.DateTime `json:"updated_at" bson:"updated_at"`
	Orders           []string          `json:"orders" bson:"orders"`
	Cart             []string          `json:"cart" bson:"cart"`
	IsActive         bool              `json:"is_active" bson:"is_active"`
	Roles            []string          `json:"roles" bson:"roles"`
	LastLogin        primitive.DateTime `json:"last_login,omitempty" bson:"last_login,omitempty"`
}


func NewUser() *User {
	currentTime := primitive.NewDateTimeFromTime(time.Now())
	return &User{
		ID:               primitive.NewObjectID(),
		CreatedAt:        currentTime,
		UpdatedAt:        currentTime,
		IsActive:         true,
		Roles:            []string{"user"},
		Orders:           []string{},
		Cart:             []string{},

	}
}

