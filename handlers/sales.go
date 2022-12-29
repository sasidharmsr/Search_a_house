package handlers

import (
	"Task/db"
	"Task/model"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// creating a user
func PostSale(dbName string) http.HandlerFunc {
	var Sales = db.GetCollection(db.DB, dbName, Collecion_sales_name)
	var Homes = db.GetCollection(db.DB, dbName, Collecion_house_name)
	return func(w http.ResponseWriter, r *http.Request) {
		var sale DummySales
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := json.NewDecoder(r.Body).Decode(&sale); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := model.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		date, err := time.Parse("2006-01-02", sale.Date) // changing date format
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := model.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}
		var house model.House
		if err = Homes.FindOne(ctx, bson.M{"_id": sale.Property}).Decode(&house); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := model.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		if house.Owner != sale.Seller {
			w.WriteHeader(http.StatusConflict)
			response := model.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Property is Under different User"}}
			json.NewEncoder(w).Encode(response)
			return
		}
		if sale.Seller == sale.Buyer {
			w.WriteHeader(http.StatusConflict)
			response := model.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "Seller and Buyer must be different Users"}}
			json.NewEncoder(w).Encode(response)
			return
		}

		newSale := model.Sales{
			ID:            primitive.NewObjectID(),
			Buyer:         sale.Buyer,
			Seller:        sale.Seller,
			Date:          date,
			Property:      sale.Property,
			TransactionId: sale.TransactionId,
		}

		res, err := Sales.InsertOne(ctx, newSale) // Inserting into Database

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := model.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}
		house.Owner = sale.Buyer
		_, err = Homes.UpdateOne(ctx, bson.M{"_id": sale.Property}, bson.M{"$set": house})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := model.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusCreated)
		response := model.Response{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": res}}
		json.NewEncoder(w).Encode(response)

	}
}

// Showing Sales by Single User
func GetSingleUserSales(dbName string) http.HandlerFunc {
	var Sales = db.GetCollection(db.DB, dbName, Collecion_sales_name)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := path.Base(fmt.Sprint(r.URL))
		var SalesResp []SalesResponse
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)

		match := bson.D{{Key: "$match", Value: bson.D{{Key: "buyer", Value: objId}}}}
		lookup1 := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "users"}, {Key: "localField", Value: "seller"}, {Key: "foreignField", Value: "_id"}, {Key: "as", Value: "seller"}}}}
		lookup2 := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "house"}, {Key: "localField", Value: "property"}, {Key: "foreignField", Value: "_id"}, {Key: "as", Value: "property"}}}}
		unwindStage1 := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$seller"}}}}
		unwindStage2 := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$property"}}}}
		pipe := mongo.Pipeline{match, lookup1, lookup2, unwindStage1, unwindStage2}

		ResultCursor, _ := Sales.Aggregate(ctx, pipe)

		if err := ResultCursor.All(ctx, &SalesResp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := model.Response{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusOK)
		response := model.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": SalesResp}}
		json.NewEncoder(w).Encode(response)
	}
}

// Delete a User from DB
func DeleteSale(dbName string) http.HandlerFunc {
	var Sales = db.GetCollection(db.DB, dbName, Collecion_sales_name)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		userId := params["id"]
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)

		result, err := Sales.DeleteOne(ctx, bson.M{"_id": objId}) // Deleting User from DB

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := model.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		// If user is not  Found in db
		if result.DeletedCount == 0 {
			w.WriteHeader(http.StatusNotFound)
			response := model.Response{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "User with specified ID not found!"}}
			json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusOK)
		response := model.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "User successfully deleted!"}}
		json.NewEncoder(w).Encode(response)
	}
}

type DummySales struct {
	Buyer         primitive.ObjectID `bson:"buyer" json:"buyer"`
	Seller        primitive.ObjectID `bson:"seller" json:"seller"`
	Property      primitive.ObjectID `bson:"property" json:"property"`
	Date          string             `json:"Date"`
	TransactionId string             `json:"transaction_id"`
}

type SalesResponse struct {
	Buyer         primitive.ObjectID `bson:"buyer" json:"buyer"`
	Seller        model.User         `bson:"seller" json:"seller"`
	Property      model.House        `bson:"property" json:"property"`
	Date          time.Time          `json:"Date"`
	TransactionId string             `json:"transaction_id"`
}
