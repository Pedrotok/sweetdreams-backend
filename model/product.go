package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Price       float64            `json:"price" bson:"price"`
	Description string             `json:"description" bson:"description"`
}
