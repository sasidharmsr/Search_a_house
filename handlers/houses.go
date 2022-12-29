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

var Collecion_house_name string = "house"
var Collecion_user_name string = "users"
var Collecion_sales_name string = "sales"

// creating a house
func PostHouse(dbName string) http.HandlerFunc {
	var Houses = db.GetCollection(db.DB, dbName, Collecion_house_name)
	return func(w http.ResponseWriter, r *http.Request) {
		var house model.House
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := json.NewDecoder(r.Body).Decode(&house); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := model.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}
		newHouse := model.House{
			Name:        house.Name,
			Address:     house.Address,
			Description: house.Description,
			Facing:      house.Facing,
			Owner:       house.Owner,
			Rate:        house.Rate,
			ID:          primitive.NewObjectID(),
			Category:    house.Category,
		}

		res, err := Houses.InsertOne(ctx, newHouse) // Inserting into Database
		if err != nil {
			fmt.Println(err)
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

// Showing All houses in DB
func GetAllHouses(dbName string) http.HandlerFunc {
	var Houses = db.GetCollection(db.DB, dbName, Collecion_house_name)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		lookup := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "users"}, {Key: "localField", Value: "owner"}, {Key: "foreignField", Value: "_id"}, {Key: "as", Value: "owner"}}}}
		unwindStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$owner"}}}}

		pipe := mongo.Pipeline{lookup, unwindStage}
		ResultCursor, _ := Houses.Aggregate(ctx, pipe)
		var ResHome []House
		if err := ResultCursor.All(ctx, &ResHome); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := model.Response{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusOK)
		response := model.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": ResHome}}
		json.NewEncoder(w).Encode(response)
	}
}

// Showing Single house
func GetSingleHouse(dbName string) http.HandlerFunc {
	var Houses = db.GetCollection(db.DB, dbName, Collecion_house_name)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		Id := path.Base(fmt.Sprint(r.URL))

		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(Id)
		match := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: objId}}}}
		lookup := bson.D{{Key: "$lookup", Value: bson.D{{Key: "from", Value: "users"}, {Key: "localField", Value: "owner"}, {Key: "foreignField", Value: "_id"}, {Key: "as", Value: "owner"}}}}
		unwindStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$owner"}}}}

		pipe := mongo.Pipeline{match, lookup, unwindStage}
		ResultCursor, _ := Houses.Aggregate(ctx, pipe)
		var ResHome []House

		if err := ResultCursor.All(ctx, &ResHome); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := model.Response{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		if (ResHome) == nil {
			w.WriteHeader(http.StatusNotFound)
			response := model.Response{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "User Not Found"}}
			json.NewEncoder(w).Encode(response)
			return
		}
		w.WriteHeader(http.StatusOK)
		response := model.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": ResHome}}
		json.NewEncoder(w).Encode(response)
	}
}

// Updating a house
func UpdateHouse(dbName string) http.HandlerFunc {
	var Houses = db.GetCollection(db.DB, dbName, Collecion_house_name)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		Id := path.Base(fmt.Sprint(r.URL))
		defer cancel()
		var home map[string]interface{}
		objId, _ := primitive.ObjectIDFromHex(Id)
		if err := json.NewDecoder(r.Body).Decode(&home); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := model.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}
		home["owner"], _ = primitive.ObjectIDFromHex(fmt.Sprint(home["owner"]))
		result, err := Houses.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": home})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := model.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		var updatedHouse model.House
		if result.MatchedCount == 1 {
			err := Houses.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedHouse)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				response := model.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				json.NewEncoder(w).Encode(response)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		response := model.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedHouse}}
		json.NewEncoder(w).Encode(response)
	}
}

// Delete a house from DB
func DeleteHouse(dbName string) http.HandlerFunc {
	var Houses = db.GetCollection(db.DB, dbName, Collecion_house_name)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		Id := params["id"]
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(Id)

		result, err := Houses.DeleteOne(ctx, bson.M{"_id": objId}) // Deleting house from DB

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := model.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		// If house is not  Found in db
		if result.DeletedCount == 0 {
			w.WriteHeader(http.StatusNotFound)
			response := model.Response{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "house with specified ID not found!"}}
			json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusOK)
		response := model.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "house successfully deleted!"}}
		json.NewEncoder(w).Encode(response)
	}
}

type House struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Address     model.Address      `json:"address"  bson:"address"`
	Owner       model.User         `json:"owner"  bson:"owner"`
	Rate        string             `json:"rate" bson:"rate"`
	Facing      string             `json:"facing"`
	Description string             `json:"description"`
	Category    string             `json:"category"`
}
