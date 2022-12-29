package handlers

import (
	"net/http"
)

func GetHousesNearBy(dbName string) http.HandlerFunc {
	//var Houses = db.GetCollection(db.DB, dbName, Collecion_house_name)
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func GetHousesbyQuery(dbName string) http.HandlerFunc {
	//var Houses = db.GetCollection(db.DB, dbName, Collecion_house_name)
	return func(w http.ResponseWriter, r *http.Request) {
		// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		// defer cancel()
		// parms:=r.URL.Query()
		// rentLessthan:=parms.Get("rentL");
		// rentGreaterthan:=parms.Get("rentG");
		// facing:=parms.Get("facing")
		// category:=parms.Get("category")

	}
}

func GetUsersbyQuery(dbName string) http.HandlerFunc {
	//var Houses = db.GetCollection(db.DB, dbName, Collecion_house_name)
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
