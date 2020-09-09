package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	Email    string             `json:"email" bson:"email"`
	Password []byte             `json:"password" bson:"password"`
}
