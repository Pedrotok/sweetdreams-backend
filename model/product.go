package model

import (
	"context"
	"log"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Product struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Price       float64            `json:"price" bson:"price"`
	Description string             `json:"description" bson:"description"`
}

func CreateProduct(name string, price float64, description string, db *mongo.Database) (*Product, error) {
	product := &Product{
		Name:        name,
		Price:       price,
		Description: description,
	}

	result, err := db.Collection("Product").InsertOne(context.TODO(), product)
	if err != nil {
		switch err.(type) {
		case mongo.WriteException:
			return nil, errors.Wrap(err, "Error while inserting data")
		default:
			return nil, errors.Wrap(err, "Error")
		}
	}
	product.ID = result.InsertedID.(primitive.ObjectID)

	return product, nil
}

func SelectProductById(id primitive.ObjectID, db *mongo.Database) (*Product, error) {
	var product Product
	err := db.Collection("Product").FindOne(context.TODO(), Product{ID: id}).Decode(&product)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return nil, errors.Wrap(err, "Product not found")
		default:
			log.Printf("Error while decode to go struct:%v\n", err)
			return nil, errors.Wrap(err, "Error while querying db")
		}
	}
	return &product, nil
}

func SelectProducts(toSkip int64, amount int64, db *mongo.Database) ([]Product, error) {
	var productList []Product

	findOptions := options.FindOptions{
		Skip:  &toSkip,
		Limit: &amount,
		Sort: bson.M{
			"_id": -1, // -1 for descending
		},
	}

	curser, err := db.Collection("Product").Find(context.TODO(), bson.M{}, &findOptions)
	if err != nil {
		log.Printf("Error while quering collection: %v\n", err)
		return nil, errors.Wrap(err, "Error  while reading data")
	}

	err = curser.All(context.Background(), &productList)
	if err != nil {
		log.Fatalf("Error in curser: %v", err)
		return nil, errors.Wrap(err, "Error  while reading data")
	}

	return productList, nil
}

func UpdateProduct(id primitive.ObjectID, updateData map[string]interface{}, db *mongo.Database) error {
	update := bson.M{
		"$set": updateData,
	}
	result, err := db.Collection("Product").UpdateOne(context.Background(), Product{ID: id}, update)
	if err != nil {
		log.Printf("Error while updating document: %v", err)
		return errors.Wrap(err, "Error while updating document")
	}
	if result.MatchedCount != 1 {
		return errors.Wrap(err, "product not found")
	}
	return nil
}
