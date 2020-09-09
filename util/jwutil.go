package util

import (
	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetToken(ID primitive.ObjectID) (string, error) {
	signingKey := []byte("keymaker")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID": ID,
	})
	tokenString, err := token.SignedString(signingKey)
	return tokenString, err
}

func VerifyToken(tokenString string) (jwt.Claims, error) {
	signingKey := []byte("keymaker")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims, err
}
