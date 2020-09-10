package controller

import (
	controller "SweetDreams/controller/requestModel"
	"SweetDreams/model"
	"SweetDreams/util"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterUser(db *mongo.Database, res http.ResponseWriter, req *http.Request) {
	registerRequest := new(controller.RegisterRequest)
	err := json.NewDecoder(req.Body).Decode(registerRequest)

	if err != nil {
		ResponseWriter(res, http.StatusBadRequest, "body json request have issues!!!", nil)
		return
	}

	_, err = model.CreateUser(registerRequest.Email, registerRequest.Password, db)

	if err != nil {
		ResponseWriter(res, http.StatusBadRequest, "password not acceptable", nil)
		return
	}

	ResponseWriter(res, http.StatusCreated, "User created", nil)
}

func Authenticate(db *mongo.Database, res http.ResponseWriter, req *http.Request) {
	registerRequest := new(controller.RegisterRequest)
	err := json.NewDecoder(req.Body).Decode(registerRequest)

	if err != nil {
		ResponseWriter(res, http.StatusNotAcceptable, "body json request have issues!!!", nil)
		return
	}

	user, err := model.AuthenticateUser(registerRequest.Email, registerRequest.Password, db)

	if err != nil {
		ResponseWriter(res, http.StatusOK, "user not found", nil)
	}

	token, err := util.GetToken(user.ID)
	res.Header().Set("Authorization", "Bearer "+token)
	ResponseWriter(res, http.StatusOK, "", token)
}
