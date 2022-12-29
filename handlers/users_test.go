package handlers

import (
	"Task/model"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Creaing a request and Response
func createNewRequestNewResponce(method, endpoint string, body io.Reader) (*http.Request, *httptest.ResponseRecorder) {
	req, _ := http.NewRequest(method, endpoint, body)
	rr := httptest.NewRecorder()
	return req, rr
}

func TestPostUser(t *testing.T) {
	// creating a dummy user
	person, _ := json.Marshal(DummyUser{Name: "king", DOB: "2003-01-02", Email: "fake@gmail.com"})

	httpHandler := http.HandlerFunc(PostUser("go-test")) // Here i used Testing DB Name
	req, rr := createNewRequestNewResponce("POST", "/person", bytes.NewBuffer(person))
	httpHandler.ServeHTTP(rr, req)

	var user model.Response
	var idObj idstuct
	_ = json.NewDecoder(rr.Body).Decode(&user)
	id, _ := json.Marshal(user.Data["data"])
	_ = json.NewDecoder(bytes.NewBuffer(id)).Decode(&idObj)
	UserId := idObj.InsertedID.Hex()
	fmt.Println(UserId)
	assert.Equal(t, rr.Code, http.StatusCreated) // checking id length
	assert.Equal(t, 24, len(UserId))

	// Bad Request
	req, rr = createNewRequestNewResponce("POST", "/person", bytes.NewBuffer([]byte("{'nam':'hat' 'addresss':'i dont know'}"))) // wrong json body
	httpHandler.ServeHTTP(rr, req)
	assert.Equal(t, rr.Code, http.StatusBadRequest)
}

func TestGetALLUser(t *testing.T) {

	httpHandler := http.HandlerFunc(GetAllUsers("go-test"))
	req, rr := createNewRequestNewResponce("GET", "/user", nil)
	httpHandler.ServeHTTP(rr, req)
	var res ResModel
	_ = json.NewDecoder(rr.Body).Decode(&res)
	assert.Equal(t, rr.Code, http.StatusOK) // checking response
}

func TestGetSingleUser(t *testing.T) {

	httpHandler := http.HandlerFunc(GetAllUsers("go-test"))
	req, rr := createNewRequestNewResponce("GET", "/user", nil)
	httpHandler.ServeHTTP(rr, req)
	var res ResModel
	_ = json.NewDecoder(rr.Body).Decode(&res)
	UserId := res.Data["data"][0].ID.Hex()

	httpHandler = http.HandlerFunc(GetSingleUser("go-test"))
	endpoint := "/user/" + UserId
	req, rr = createNewRequestNewResponce("GET", endpoint, nil)
	httpHandler.ServeHTTP(rr, req)

	var user model.Response
	var idObj model.User
	_ = json.NewDecoder(rr.Body).Decode(&user)
	id, _ := json.Marshal(user.Data["data"])
	_ = json.NewDecoder(bytes.NewBuffer(id)).Decode(&idObj)
	ResponceUserId := idObj.ID.Hex()

	assert.Equal(t, UserId, ResponceUserId) //checking USERID
	assert.Equal(t, rr.Code, http.StatusOK)

}

func TestUpdateUser(t *testing.T) {

	httpHandler := http.HandlerFunc(GetAllUsers("go-test"))
	req, rr := createNewRequestNewResponce("GET", "/user", nil)
	httpHandler.ServeHTTP(rr, req)
	var res ResModel
	_ = json.NewDecoder(rr.Body).Decode(&res)
	UserId := res.Data["data"][0].ID.Hex()

	person, _ := json.Marshal(map[string]string{"name": "raju"}) // Updating User
	UserName := "raju"

	httpHandler = http.HandlerFunc(UpdateUser("go-test"))
	endpoint := "/user/" + UserId
	req, rr = createNewRequestNewResponce("PUT", endpoint, bytes.NewBuffer(person))
	httpHandler.ServeHTTP(rr, req)

	var user model.Response
	var idObj model.User
	_ = json.NewDecoder(rr.Body).Decode(&user)
	id, _ := json.Marshal(user.Data["data"])
	_ = json.NewDecoder(bytes.NewBuffer(id)).Decode(&idObj)
	ResponceUserId := idObj.ID.Hex()

	assert.Equal(t, UserId, ResponceUserId) //Checking User Id
	assert.Equal(t, UserName, idObj.Name)   // Checking Updated UserName
	assert.Equal(t, rr.Code, http.StatusOK)
}

type idstuct struct {
	InsertedID primitive.ObjectID
}

type ResModel struct {
	Status  int                     `json:"status"`
	Message string                  `json:"message"`
	Data    map[string][]model.User `json:"data"`
}
