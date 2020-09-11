package controller

import (
	controller "SweetDreams/controller/requestModel"
	"SweetDreams/model"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// results count per page
var limit int64 = 10

func CreateProduct(db *mongo.Database, res http.ResponseWriter, req *http.Request) error {
	request := new(controller.CreateProductRequest)
	err := json.NewDecoder(req.Body).Decode(request)
	if err != nil {
		return StatusError{http.StatusBadRequest, errors.Wrap(err, "Failed to decode request")}
	}

	_, err = model.CreateProduct(request.Name, request.Price, request.Description, db)

	if err != nil {
		return err
	}

	return ResponseWriter(res, http.StatusCreated, "", nil)
}

func GetProduct(db *mongo.Database, res http.ResponseWriter, req *http.Request) error {
	var params = mux.Vars(req)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		return StatusError{http.StatusBadRequest, errors.Wrap(err, "Bad request")}
	}

	product, err := model.SelectProductById(id, db)

	if err != nil {
		return StatusError{http.StatusNotFound, errors.Wrap(err, "Product not found")}
	}

	return ResponseWriter(res, http.StatusOK, "", product)
}

func GetAllProducts(db *mongo.Database, res http.ResponseWriter, req *http.Request) error {
	pageString := req.FormValue("page")
	page, err := strconv.ParseInt(pageString, 10, 64)
	if err != nil {
		page = 0
	}

	productList, err := model.SelectProducts(page*limit, limit, db)

	if err != nil {
		return StatusError{http.StatusNotFound, errors.Wrap(err, "Can't query products")}
	}

	return ResponseWriter(res, http.StatusOK, "", productList)
}

func UpdateProduct(db *mongo.Database, res http.ResponseWriter, req *http.Request) error {
	var updateData map[string]interface{}
	err := json.NewDecoder(req.Body).Decode(&updateData)
	if err != nil {
		return StatusError{http.StatusBadRequest, errors.Wrap(err, "Bad request")}
	}

	var params = mux.Vars(req)
	oid, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		return StatusError{http.StatusBadRequest, errors.Wrap(err, "Bad request")}
	}

	err = model.UpdateProduct(oid, updateData, db)

	if err != nil {
		return StatusError{http.StatusNotFound, errors.Wrap(err, "Error updating product")}
	}

	return ResponseWriter(res, http.StatusOK, "", nil)
}
