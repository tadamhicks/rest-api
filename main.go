package main

import (
	"encoding/json"
	"log"
	"net/http"
	"gopkg.in/mgo.v2/bson"
	"github.com/gorilla/mux"
	"github.com/dgrijalva/jwt-go"
	"time"
	"github.com/auth0/go-jwt-middleware"
	. "github.com/tadamhicks/rest-api/config"
	. "github.com/tadamhicks/rest-api/dao"
	. "github.com/tadamhicks/rest-api/models"
)


var config = Config{}
var dao = PersonDAO{}
var mySigningKey = []byte("secret")


var GetToken = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["admin"] = true
	claims["name"] = "Pepe LePeux"
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, _ := token.SignedString(mySigningKey)
	w.Write([]byte(tokenString))
})


var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
	  return mySigningKey, nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})


var GetPeople = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	person, err := dao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, person)
})



var UpdatePerson = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	defer r.Body.Close()
	var person Person
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Update(params["id"], person); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
})


var GetPerson = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	person, err := dao.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Person ID")
		return
	}
	respondWithJson(w, http.StatusOK, person)
})

var CreatePerson = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var person Person
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	person.ID = bson.NewObjectId()
	if err := dao.Insert(person); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, person)
})


var DeletePerson = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	err := dao.Delete(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Person ID")
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
})


func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}


func init() {
	config.Read()
	dao.Server = config.Server
	dao.Database = config.Database
	dao.Username = config.Username
	dao.Password = config.Password
	dao.Connect()
}

func main() {
	router := mux.NewRouter()
	router.Handle("/get-token", GetToken).Methods("GET")
	router.Handle("/people", GetPeople).Methods("GET")
	router.Handle("/people/{id}", jwtMiddleware.Handler(UpdatePerson)).Methods("PUT")
	router.Handle("/people/{id}", GetPerson).Methods("GET")
	router.Handle("/people", jwtMiddleware.Handler(CreatePerson)).Methods("POST")
	router.Handle("/people/{id}", jwtMiddleware.Handler(DeletePerson)).Methods("DELETE")
	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Fatal(err)
	}
}
