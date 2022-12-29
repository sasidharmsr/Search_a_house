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
)

// creating a user
func PostUser(dbName string) http.HandlerFunc {
	var Users = db.GetCollection(db.DB, dbName, Collecion_user_name)
	return func(w http.ResponseWriter, r *http.Request) {
		var user DummyUser
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := model.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		date, err := time.Parse("2006-01-02", user.DOB) // changing date format
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := model.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}
		newUser := model.User{
			Name:      user.Name,
			DOB:       date,
			PhoneNo:   user.PhoneNo,
			Email:     user.Email,
			UserName:  user.UserName,
			ID:        primitive.NewObjectID(),
			CreatedAt: primitive.NewDateTimeFromTime(time.Now().UTC()),
		}

		res, err := Users.InsertOne(ctx, newUser) // Inserting into Database
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

// Showing All users in DB
func GetAllUsers(dbName string) http.HandlerFunc {
	var Users = db.GetCollection(db.DB, dbName, Collecion_user_name)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var users []model.User
		results, err := Users.Find(ctx, bson.M{}) //Here we will get all users from DB

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := model.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		//reading users from the db
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleUser model.User
			if err = results.Decode(&singleUser); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				response := model.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				json.NewEncoder(w).Encode(response)
			}
			users = append(users, singleUser)
		}

		w.WriteHeader(http.StatusOK)
		response := model.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": users}}
		json.NewEncoder(w).Encode(response)
	}
}

// Showing Single User
func GetSingleUser(dbName string) http.HandlerFunc {
	var Users = db.GetCollection(db.DB, dbName, Collecion_user_name)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := path.Base(fmt.Sprint(r.URL))

		defer cancel()
		var user model.User
		objId, _ := primitive.ObjectIDFromHex(userId)
		err := Users.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := model.Response{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}
		w.WriteHeader(http.StatusOK)
		response := model.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": user}}
		json.NewEncoder(w).Encode(response)
	}
}

// Updating a User
func UpdateUser(dbName string) http.HandlerFunc {
	var Users = db.GetCollection(db.DB, dbName, Collecion_user_name)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := path.Base(fmt.Sprint(r.URL))
		defer cancel()
		var user map[string]interface{}
		objId, _ := primitive.ObjectIDFromHex(userId)
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response := model.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		if user["dob"] != nil {
			user["dob"], _ = time.Parse("2006-01-02", fmt.Sprintf("%f", user["dob"]))
		}

		result, err := Users.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": user})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response := model.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			json.NewEncoder(w).Encode(response)
			return
		}

		var updatedUser model.User
		if result.MatchedCount == 1 {
			err := Users.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedUser)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				response := model.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				json.NewEncoder(w).Encode(response)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		response := model.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedUser}}
		json.NewEncoder(w).Encode(response)
	}
}

// Delete a User from DB
func DeleteUser(dbName string) http.HandlerFunc {
	var Users = db.GetCollection(db.DB, dbName, Collecion_user_name)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		params := mux.Vars(r)
		userId := params["id"]
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)

		result, err := Users.DeleteOne(ctx, bson.M{"_id": objId}) // Deleting User from DB

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

// Showing All USERS Near By A Particlur USer
// func GetNearByUsers(dbName string) http.HandlerFunc {
// 	var Users = db.GetCollection(db.DB, dbName)
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 		defer cancel()
// 		var users []NearestUsersResp
// 		results, err := Users.Find(ctx, bson.M{}) //Here we will get all users from DB
// 		params := mux.Vars(r)
// 		userId := params["id"]

// 		if err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			response := model.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
// 			json.NewEncoder(w).Encode(response)
// 			return
// 		}

// 		//reading a user from the db
// 		var user model.User
// 		objId, _ := primitive.ObjectIDFromHex(userId)
// 		err = Users.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
// 		if err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			response := model.Response{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
// 			json.NewEncoder(w).Encode(response)
// 			return
// 		}

// 		defer results.Close(ctx)
// 		for results.Next(ctx) {
// 			var singleUser model.User
// 			var singleresp NearestUsersResp
// 			if err = results.Decode(&singleUser); err != nil {
// 				w.WriteHeader(http.StatusInternalServerError)
// 				response := model.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
// 				json.NewEncoder(w).Encode(response)
// 			}
// 			if singleUser.ID.Hex() != userId {
// 				singleresp.UserId = singleUser.ID.Hex()
// 				singleresp.UserName = singleUser.Name
// 				singleresp.Address = singleUser.Address.Addresss
// 				singleresp.Distance = DistancebetweenLoc(user.Address.Latitude, user.Address.Longitude, singleUser.Address.Latitude, singleUser.Address.Longitude, "K")
// 				users = append(users, singleresp)
// 			}
// 		}
// 		sort.Slice(users, func(i, j int) bool { return users[i].Distance < users[j].Distance })
// 		w.WriteHeader(http.StatusOK)
// 		response := model.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": users}}
// 		json.NewEncoder(w).Encode(response)
// 	}
// }

// This stuct is only for changing DOB data type
type DummyUser struct {
	UserName  string             `json:"user_name"`
	Email     string             `json:"email"`
	Name      string             `json:"name"`
	PhoneNo   string             `json:"phone_number"`
	ID        primitive.ObjectID `bson:"_id"`
	DOB       string             `json:"dob"`
	CreatedAt primitive.DateTime `json:"created_at"`
}
