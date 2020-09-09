package controller

import (
	"SweetDreams/model"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// results count per page
var limit int64 = 10

func CreateProduct(db *mongo.Database, res http.ResponseWriter, req *http.Request) {
	product := new(model.Product)
	err := json.NewDecoder(req.Body).Decode(product)
	if err != nil {
		ResponseWriter(res, http.StatusBadRequest, "body json request have issues!!!", nil)
		return
	}
	result, err := db.Collection("Product").InsertOne(nil, product)
	if err != nil {
		switch err.(type) {
		case mongo.WriteException:
			ResponseWriter(res, http.StatusNotAcceptable, "Error while inserting data.", nil)
		default:
			ResponseWriter(res, http.StatusInternalServerError, "Error while inserting data.", nil)
		}
		return
	}
	product.ID = result.InsertedID.(primitive.ObjectID)
	ResponseWriter(res, http.StatusCreated, "", product)
}

func GetProduct(db *mongo.Database, res http.ResponseWriter, req *http.Request) {
	var params = mux.Vars(req)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		ResponseWriter(res, http.StatusBadRequest, "id that you sent is wrong!!!", nil)
		return
	}
	var product model.Product
	err = db.Collection("Product").FindOne(nil, model.Product{ID: id}).Decode(&product)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			ResponseWriter(res, http.StatusNotFound, "product not found", nil)
		default:
			log.Printf("Error while decode to go struct:%v\n", err)
			ResponseWriter(res, http.StatusInternalServerError, "there is an error on server!!!", nil)
		}
		return
	}
	ResponseWriter(res, http.StatusOK, "", product)
}

func GetAllProducts(db *mongo.Database, res http.ResponseWriter, req *http.Request) {
	var productList []model.Product
	pageString := req.FormValue("page")
	page, err := strconv.ParseInt(pageString, 10, 64)
	if err != nil {
		page = 0
	}
	page = page * limit
	findOptions := options.FindOptions{
		Skip:  &page,
		Limit: &limit,
		Sort: bson.M{
			"_id": -1, // -1 for descending and 1 for ascending
		},
	}
	curser, err := db.Collection("Product").Find(nil, bson.M{}, &findOptions)
	if err != nil {
		log.Printf("Error while quering collection: %v\n", err)
		ResponseWriter(res, http.StatusInternalServerError, "Error happend while reading data", nil)
		return
	}
	err = curser.All(context.Background(), &productList)
	if err != nil {
		log.Fatalf("Error in curser: %v", err)
		ResponseWriter(res, http.StatusInternalServerError, "Error happend while reading data", nil)
		return
	}
	ResponseWriter(res, http.StatusOK, "", productList)
}

func UpdateProduct(db *mongo.Database, res http.ResponseWriter, req *http.Request) {
	var updateData map[string]interface{}
	err := json.NewDecoder(req.Body).Decode(&updateData)
	if err != nil {
		ResponseWriter(res, http.StatusBadRequest, "json body is incorrect", nil)
		return
	}
	// we dont handle the json decode return error because all our fields have the omitempty tag.
	var params = mux.Vars(req)
	oid, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		ResponseWriter(res, http.StatusBadRequest, "id that you sent is wrong!!!", nil)
		return
	}
	update := bson.M{
		"$set": updateData,
	}
	result, err := db.Collection("Product").UpdateOne(context.Background(), model.Product{ID: oid}, update)
	if err != nil {
		log.Printf("Error while updateing document: %v", err)
		ResponseWriter(res, http.StatusInternalServerError, "error in updating document!!!", nil)
		return
	}
	if result.MatchedCount == 1 {
		ResponseWriter(res, http.StatusAccepted, "", &updateData)
	} else {
		ResponseWriter(res, http.StatusNotFound, "product not found", nil)
	}
}
