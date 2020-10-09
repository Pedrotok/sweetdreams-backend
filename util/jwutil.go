package util

import (
	"errors"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TokenType int

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
}

const (
	Access TokenType = iota
	Refresh
)

func getKey(tokenType TokenType) string {
	if tokenType == Access {
		return viper.GetString("JwtAccessKey")
	}
	return viper.GetString("JwtRefreshKey")
}

func CreateToken(ID primitive.ObjectID) (*TokenDetails, error) {
	jwtAccessKey := viper.GetString("JwtAccessKey")
	accessKey := []byte(jwtAccessKey)

	jwtRefreshKey := viper.GetString("JwtRefreshKey")
	refreshKey := []byte(jwtRefreshKey)

	td := &TokenDetails{}
	atExpires := time.Now().Add(time.Minute * 15).Unix()
	rtExpires := time.Now().Add(time.Hour * 24 * 7).Unix()

	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["ID"] = ID
	atClaims["exp"] = atExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(accessKey))
	if err != nil {
		return nil, err
	}

	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["ID"] = ID
	rtClaims["exp"] = rtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString(refreshKey)
	if err != nil {
		return nil, err
	}

	return td, nil
}

func VerifyToken(tokenString string, tokenType TokenType) (*jwt.Token, error) {
	jwtKeyString := getKey(tokenType)
	accessKey := []byte(jwtKeyString)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(accessKey), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func GetToken(tokenString string, tokenType TokenType) (*jwt.Token, error) {
	token, err := VerifyToken(tokenString, tokenType)
	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, errors.New("Invalid token")
	}

	return token, nil
}
