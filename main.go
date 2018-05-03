package main

import (
	"encoding/json"
	"log"
	"net/http"
	"gopkg.in/mgo.v2/bson"
	"github.com/gorilla/mux"
	. "github.com/tadamhicks/rest-api/config"
	. "github.com/tadamhicks/rest-api/dao"
	. "github.com/tadamhicks/rest-api/models"
)


var config = Config{}
var dao = PersonDAO{}

func GetPeople(w http.ResponseWriter, r *http.Request) {
	person, err := dao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, person)
}


func UpdatePerson(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var person Person
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Update(person); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func GetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	person, err := dao.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Person ID")
		return
	}
	respondWithJson(w, http.StatusOK, person)
}

func CreatePerson(w http.ResponseWriter, r *http.Request) {
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
}


func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	err := dao.Delete(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Person ID")
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}
/*
func DeletePerson(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var person Person
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Delete(person); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}
*/


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
	dao.Connect()
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/people", GetPeople).Methods("GET")
	router.HandleFunc("/people", UpdatePerson).Methods("PUT")
	router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
	router.HandleFunc("/people", CreatePerson).Methods("POST")
	router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")
	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Fatal(err)
	}
}
