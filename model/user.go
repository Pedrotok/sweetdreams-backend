package model

import (
	"context"

	"github.com/pkg/errors"
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
		return nil, errors.Wrap(err, "password not valid")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "password not hashable")
	}

	user := new(User)
	user.Email = email
	user.Password = hash

	result, err := db.Collection("User").InsertOne(context.TODO(), user)
	if err != nil {
		switch err.(type) {
		case mongo.WriteException:
			return nil, errors.Wrap(err, "Error while inserting data")
		default:
			return nil, errors.Wrap(err, "Error")
		}
	}
	user.ID = result.InsertedID.(primitive.ObjectID)

	return user, nil
}

func AuthenticateUser(email string, password string, db *mongo.Database) (*User, error) {
	user := new(User)
	err := db.Collection("User").FindOne(context.TODO(), bson.D{{Key: "email", Value: email}}).Decode(user)
	if err != nil {
		return nil, errors.Wrap(err, "User with email not found")
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		return nil, errors.Wrap(err, "Password not matching")
	}

	return user, nil
}

func SelectUserById(id primitive.ObjectID, db *mongo.Database) (*User, error) {
	user := new(User)
	err := db.Collection("User").FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(user)
	if err != nil {
		return nil, errors.Wrap(err, "User not found")
	}

	return user, nil
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password too short")
	}
	return nil
}
