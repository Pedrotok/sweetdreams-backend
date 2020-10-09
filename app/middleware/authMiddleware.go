package app

import (
	"SweetDreams/model"
	"SweetDreams/util"

	"net/http"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AuthMiddleware(db *mongo.Database, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tokenString := util.ExtractTokenString(r)
		if len(tokenString) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Missing Authorization Header"))
			return
		}

		token, err := util.GetToken(tokenString, util.Access)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Error verifying JWT token: " + err.Error()))
			return
		}

		userId := token.Claims.(jwt.MapClaims)["ID"].(string)

		oid, err := primitive.ObjectIDFromHex(userId)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Problem parsing user id"))
			return
		}

		_, err = model.SelectUserById(oid, db)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("User not found"))
			return
		}
		r.Header.Set("user-id", userId)

		next.ServeHTTP(w, r)
	}
}
