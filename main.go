package main

import (
	"Task/db"
	"Task/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter() // Intiate a Router

	//run database
	db.ConnectDB()

	dbName := "Search_A_House"

	// Routes for user
	router.HandleFunc("/user", handlers.PostUser(dbName)).Methods("POST")
	router.HandleFunc("/user", handlers.GetAllUsers(dbName)).Methods("GET")
	router.HandleFunc("/user/{id}", handlers.GetSingleUser(dbName)).Methods("GET")
	router.HandleFunc("/user/{id}", handlers.UpdateUser(dbName)).Methods("PUT")
	router.HandleFunc("/user/{id}", handlers.DeleteUser(dbName)).Methods("DELETE")

	router.HandleFunc("/home", handlers.PostHouse(dbName)).Methods("POST")
	router.HandleFunc("/home", handlers.GetAllHouses(dbName)).Methods("GET")
	router.HandleFunc("/home/{id}", handlers.GetSingleHouse(dbName)).Methods("GET")
	router.HandleFunc("/home/{id}", handlers.UpdateHouse(dbName)).Methods("PUT")
	router.HandleFunc("/home/{id}", handlers.DeleteHouse(dbName)).Methods("DELETE")

	router.HandleFunc("/sales", handlers.PostSale(dbName)).Methods("POST")
	router.HandleFunc("/sales/{id}", handlers.DeleteSale(dbName)).Methods("DELETE")
	router.HandleFunc("/sales/{id}", handlers.GetSingleUserSales(dbName)).Methods("GET")

	router.HandleFunc("/filterhome", handlers.GetHousesbyQuery(dbName)).Methods("GET")
	router.HandleFunc("/nearbyhome", handlers.GetHousesNearBy(dbName)).Methods("POST")
	router.HandleFunc("/user/homes/{id}", handlers.GetUsersbyQuery(dbName)).Methods("GET")
	http.ListenAndServe(":8090", router)
}
