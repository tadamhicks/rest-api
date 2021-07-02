package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	beeline "github.com/honeycombio/beeline-go"
	"github.com/honeycombio/beeline-go/wrappers/hnynethttp"
	"github.com/kelseyhightower/envconfig"
	"github.com/tadamhicks/rest-api/dao"
	"github.com/tadamhicks/rest-api/models"
	"gopkg.in/mgo.v2/bson"
)

//var config = Config{}
var pao = dao.PersonDAO{}
var mySigningKey = []byte("secret")

type Config struct {
	Server   string `default:"127.0.0.1"`
	Port     string `default:"27017"`
	Database string `required:"true"`
	Username string `required:"true"`
	Password string `required:"true"`
}

/*
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
*/

var GetPeople = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	beeline.AddField(r.Context(), "email", "one@two.com")
	person, err := pao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, person)
})

var UpdatePerson = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	defer r.Body.Close()
	var person models.Person
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := pao.Update(params["id"], person); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
})

var GetPerson = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	person, err := pao.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Person ID")
		return
	}
	respondWithJson(w, http.StatusOK, person)
})

var CreatePerson = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var person models.Person
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	person.ID = bson.NewObjectId()
	if err := pao.Insert(person); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, person)
})

var DeletePerson = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	err := pao.Delete(params["id"])
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
	//config.Read()
	//fmt.Println("CONFIG:\n")
	//fmt.Printf("%+v\n", config)
	var c Config
	err := envconfig.Process("mongo", &c)
	if err != nil {
		log.Fatal(err.Error())
		//log.Fatalf("Failed to parse ENV")
	}
	output := strings.Join([]string{c.Server, c.Port}, ":")
	pao.Server = output
	pao.Database = c.Database
	pao.Username = c.Username
	pao.Password = c.Password
	pao.Connect()

	beeline.Init(beeline.Config{
		WriteKey:    c.Honeycombkey,
		Dataset:     c.Honeycombdataset,
		ServiceName: c.Servicename,
		STDOUT:      true,
	})
}

func main() {
	beeline.Init(beeline.Config{
		WriteKey:    "c4b05d6b2259d9d6fca768d4ba9c811a",
		Dataset:     "restful-sleep",
		ServiceName: "restful-sleep-svc",
		STDOUT:      true,
	})
	router := mux.NewRouter()
	//router.Handle("/get-token", GetToken).Methods("GET")
	router.Handle("/people", GetPeople).Methods("GET")
	router.Handle("/people/{id}", UpdatePerson).Methods("PUT")
	router.Handle("/people/{id}", GetPerson).Methods("GET")
	router.Handle("/people", CreatePerson).Methods("POST")
	router.Handle("/people/{id}", DeletePerson).Methods("DELETE")
	if err := http.ListenAndServe(":8000", hnynethttp.WrapHandler(router)); err != nil {
		log.Fatal(err)
	}
}
