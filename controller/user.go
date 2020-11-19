package controller

import (
	controller "SweetDreams/controller/requestModel"
	"SweetDreams/model"
	"SweetDreams/util"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterUser(db *mongo.Database, res http.ResponseWriter, req *http.Request) error {
	registerRequest := new(controller.RegisterRequest)
	err := json.NewDecoder(req.Body).Decode(registerRequest)

	if err != nil {
		return StatusError{http.StatusBadRequest, errors.Wrap(err, "Failed to decode request\n")}
	}

	_, err = model.CreateUser(registerRequest.Email, registerRequest.Password, db)

	if err != nil {
		return StatusError{http.StatusBadRequest, errors.Wrap(err, "Couldn't create user\n")}
	}

	return ResponseWriter(res, http.StatusCreated, "User created\n", nil)
}

func Authenticate(db *mongo.Database, res http.ResponseWriter, req *http.Request) error {
	registerRequest := new(controller.RegisterRequest)
	err := json.NewDecoder(req.Body).Decode(registerRequest)

	if err != nil {
		return StatusError{http.StatusBadRequest, errors.Wrap(err, "Failed to decode request\n")}
	}

	user, err := model.AuthenticateUser(registerRequest.Email, registerRequest.Password, db)

	if err != nil {
		return StatusError{http.StatusUnauthorized, errors.Wrap(err, "Couldn't authenticate user\n")}
	}

	ts, err := util.CreateToken(user.ID)

	if err != nil {
		return StatusError{http.StatusUnauthorized, errors.Wrap(err, "Couldn't create token\n")}
	}

	response := map[string]interface{}{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
		"user_data":     user,
	}

	return ResponseWriter(res, http.StatusCreated, "", response)
}
