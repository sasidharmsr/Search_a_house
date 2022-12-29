package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// My User Model
type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	UserName  string             `json:"user_name"`
	Email     string             `json:"email"`
	Name      string             `json:"name"`
	PhoneNo   string             `json:"phone_number"`
	DOB       time.Time          `json:"dob"`
	CreatedAt primitive.DateTime `json:"created_at"`
}

type House struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Address     Address            `json:"address"  bson:"address"`
	Owner       primitive.ObjectID `json:"owner"  bson:"owner"`
	Rate        string             `json:"rate" bson:"rate"`
	Facing      string             `json:"facing"`
	Description string             `json:"description"`
	Category    string             `json:"category"`
}

type Sales struct {
	ID            primitive.ObjectID `bson:"_id"`
	Buyer         primitive.ObjectID `bson:"buyer" json:"buyer"`
	Seller        primitive.ObjectID `bson:"seller" json:"seller"`
	Property      primitive.ObjectID `bson:"property" json:"property"`
	Date          time.Time          `json:"Date"`
	TransactionId string             `json:"transaction_id"`
}

// Model For Returning Response
type Response struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type Address struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Addresss  string  `json:"address"`
}
