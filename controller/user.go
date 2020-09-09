package controller

import (
	controller "SweetDreams/controller/model"
	"SweetDreams/model"
	"SweetDreams/util"
	"context"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(db *mongo.Database, res http.ResponseWriter, req *http.Request) {
	user := new(model.User)
	err := json.NewDecoder(req.Body).Decode(user)

	if err != nil {
		ResponseWriter(res, http.StatusBadRequest, "body json request have issues!!!", nil)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ResponseWriter(res, http.StatusBadRequest, "password not acceptable", nil)
		return
	}
	user.Password = hash

	result, err := db.Collection("User").InsertOne(nil, user)
	if err != nil {
		switch err.(type) {
		case mongo.WriteException:
			ResponseWriter(res, http.StatusNotAcceptable, "Error while inserting data.", nil)
		default:
			ResponseWriter(res, http.StatusInternalServerError, "Error while inserting data.", nil)
		}
		return
	}
	user.ID = result.InsertedID.(primitive.ObjectID)
	ResponseWriter(res, http.StatusCreated, "", user)
}

func Authenticate(db *mongo.Database, res http.ResponseWriter, req *http.Request) {
	registerRequest := new(controller.RegisterRequest)
	err := json.NewDecoder(req.Body).Decode(registerRequest)

	if err != nil {
		ResponseWriter(res, http.StatusNotAcceptable, "body json request have issues!!!", nil)
		return
	}

	// create a value into which the result can be decoded
	user := new(model.User)
	filter := bson.D{{"Email", registerRequest.Email}}
	err = db.Collection("User").FindOne(context.TODO(), filter).Decode(user)
	if err != nil {
		ResponseWriter(res, http.StatusBadRequest, "user not found", nil)
		return
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(registerRequest.Password)); err != nil {
		ResponseWriter(res, http.StatusBadRequest, "user not found", nil)
		return
	}

	token, err := util.GetToken(user.ID)
	res.Header().Set("Authorization", "Bearer "+token)
	ResponseWriter(res, http.StatusOK, "", token)
}
