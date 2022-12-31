package handlers

import (
	"Task/db"
	"Task/model"
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// You can get all houses sorted by distance from a specific location
func GetHousesNearBy(dbName string) http.HandlerFunc {
	var Houses = db.GetCollection(db.DB, dbName, Collecion_house_name)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var address model.Address
		if err := json.NewDecoder(r.Body).Decode(&address); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := model.Response{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}
		var houses []LocationFilterResp
		ResultCursor, _ := Houses.Find(ctx, bson.D{})
		for ResultCursor.Next(ctx) {
			var singleHouse model.House
			var singleresp LocationFilterResp
			if err := ResultCursor.Decode(&singleHouse); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				response := model.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				json.NewEncoder(w).Encode(response)
			}
			distance := DistancebetweenLoc(singleHouse.Address.Latitude, singleHouse.Address.Longitude, address.Latitude, address.Longitude, "K")
			val, _ := json.Marshal(singleHouse)
			_ = json.Unmarshal(val, &singleresp)
			singleresp.Distance = distance
			houses = append(houses, singleresp)
		}
		sort.Slice(houses, func(i, j int) bool { return houses[i].Distance < houses[j].Distance })
		w.WriteHeader(http.StatusOK)
		response := model.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": houses}}
		json.NewEncoder(w).Encode(response)
	}
}

// You can get filtered houses based on rate , facing ,category..
func GetHousesbyQuery(dbName string) http.HandlerFunc {
	var Houses = db.GetCollection(db.DB, dbName, Collecion_house_name)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		parms := r.URL.Query()
		rentLessthan := parms.Get("rateL")
		rentGreaterthan := parms.Get("rateG")
		facing := parms.Get("facing")
		category := parms.Get("category")
		var query []bson.D
		if rentLessthan != "" {
			query = append(query, bson.D{{Key: "rate", Value: bson.D{{Key: "$gte", Value: rentLessthan}}}})
		}
		if rentGreaterthan != "" {
			query = append(query, bson.D{{Key: "rate", Value: bson.D{{Key: "$lte", Value: rentGreaterthan}}}})
		}
		if facing != "" {
			query = append(query, bson.D{{Key: "facing", Value: facing}})
		}
		if category != "" {
			query = append(query, bson.D{{Key: "category", Value: category}})
		}
		match := bson.D{{Key: "$match", Value: bson.D{{Key: "$and", Value: query}}}}
		pipe := mongo.Pipeline{match}
		ResultCursor, err := Houses.Aggregate(ctx, pipe)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := model.Response{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}
		var ResHome []model.House
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

// get list of houses of a specific user
func GetUsersbyQuery(dbName string) http.HandlerFunc {
	Houses := db.GetCollection(db.DB, dbName, Collecion_house_name)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		params := mux.Vars(r)
		userId := params["id"]
		objId, _ := primitive.ObjectIDFromHex(userId)
		var houses []model.House
		RespCursor, _ := Houses.Find(ctx, bson.D{{Key: "owner", Value: objId}})
		if err := RespCursor.All(ctx, &houses); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := model.Response{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}
		w.WriteHeader(http.StatusOK)
		response := model.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": houses}}
		json.NewEncoder(w).Encode(response)
	}
}

type LocationFilterResp struct {
	Distance    float64            `json:"distance"`
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Address     model.Address      `json:"address"  bson:"address"`
	Rate        string             `json:"rate" bson:"rate"`
	Facing      string             `json:"facing"`
	Description string             `json:"description"`
	Category    string             `json:"category"`
}
