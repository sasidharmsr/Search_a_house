package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	url := "mongodb+srv://sasidhar431:1UlsXQ1EtlQJPWd5@cluster0.rlsiggd.mongodb.net/go?retryWrites=true&w=majority"
	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")
	return client
}

// Client instance
var DB *mongo.Client = ConnectDB()

// getting database collections of specific db
func GetCollection(client *mongo.Client, dbName string, collecion_name string) *mongo.Collection {
	collection := client.Database(dbName).Collection(collecion_name) // I created a collection with name gotest in db
	return collection
}
