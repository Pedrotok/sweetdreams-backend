package controller

import (
	"SweetDreams/util"
	"encoding/json"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func RefreshToken(db *mongo.Database, res http.ResponseWriter, req *http.Request) error {
	var refreshRequest map[string]interface{}
	err := json.NewDecoder(req.Body).Decode(&refreshRequest)
	if err != nil {
		return StatusError{http.StatusBadRequest, errors.Wrap(err, "Bad request")}
	}
	refreshToken := refreshRequest["refresh_token"]

	token, err := util.GetToken(refreshToken.(string), util.Refresh)

	if err != nil {
		return StatusError{http.StatusUnauthorized, errors.Wrap(err, "Invalid refresh token")}
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return StatusError{http.StatusUnauthorized, errors.Wrap(err, "Couldn't get claims")}
	}

	userId, ok := claims["ID"].(string)
	if !ok {
		return StatusError{http.StatusUnauthorized, errors.Wrap(err, "Couldn't get userId")}
	}

	oid, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		return StatusError{http.StatusUnauthorized, errors.Wrap(err, "Couldn't proccess userId")}
	}

	ts, err := util.CreateToken(oid)
	if err != nil {
		return StatusError{http.StatusUnauthorized, errors.Wrap(err, "Couldn't create token\n")}
	}

	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}

	return ResponseWriter(res, http.StatusCreated, "", tokens)
}
