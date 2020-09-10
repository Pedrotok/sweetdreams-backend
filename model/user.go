package model

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email    string             `json:"email" bson:"email"`
	Password []byte             `json:"password" bson:"password"`
}

func CreateUser(email string, password string, db *mongo.Database) (*User, error) {
	err := validatePassword(password)
	if err != nil {
		return nil, errors.New("password not valid")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("password not acceptable")
	}

	user := new(User)
	user.Email = email
	user.Password = hash

	result, err := db.Collection("User").InsertOne(context.TODO(), user)
	if err != nil {
		switch err.(type) {
		case mongo.WriteException:
			log.Printf(err.Error())
			return nil, errors.New("Error while inserting data")
		default:
			return nil, errors.New("Error")
		}
	}
	user.ID = result.InsertedID.(primitive.ObjectID)

	return user, nil
}

func AuthenticateUser(email string, password string, db *mongo.Database) (*User, error) {
	user := new(User)
	err := db.Collection("User").FindOne(context.TODO(), bson.D{{"email", email}}).Decode(user)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		return nil, err
	}

	return user, nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password too short")
	}
	return nil
}
