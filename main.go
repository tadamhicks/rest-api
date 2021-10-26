package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/tadamhicks/rest-api/dao"
	"github.com/tadamhicks/rest-api/models"
	"gopkg.in/mgo.v2/bson"
	"go.opentelemetry.io/otel"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	middleware "go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
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

type HoneyCfg struct {
	Apikey      string `required:"true`
	Dataset     string `required:"true"`
	Servicename string `required:"true"`
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
	//beeline.AddField(r.Context(), "email", "one@two.com")
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
	//beeline.AddField(r.Context(), "person", person)
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
})

var GetPerson = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	person, err := pao.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Person ID")
		return
	}
	//beeline.AddField(r.Context(), "person", person)
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
	//beeline.AddField(r.Context(), "person", person)
	respondWithJson(w, http.StatusCreated, person)
})

var DeletePerson = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	err := pao.Delete(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Person ID")
		return
	}
	//beeline.AddField(r.Context(), "id", params["id"])
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

}

func main() {
	/*
	var h HoneyCfg
	err := envconfig.Process("honeycomb", &h)
	if err != nil {
		log.Fatal(err.Error())
		//log.Fatalf("Failed to parse ENV")
	}
	 */

	router := mux.NewRouter()
	//router.Use(hnygorilla.Middleware)
	//router.Handle("/get-token", GetToken).Methods("GET")
	router.Handle("/people", GetPeople).Methods("GET")
	router.Handle("/people/{id}", UpdatePerson).Methods("PUT")
	router.Handle("/people/{id}", GetPerson).Methods("GET")
	router.Handle("/people", CreatePerson).Methods("POST")
	router.Handle("/people/{id}", DeletePerson).Methods("DELETE")
	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Fatal(err)
	}
}
/*
func initOtelTracing(log logrus.FieldLogger) {
	otlpendpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if otlpendpoint == "" {
		otlpendpoint = "api.honeycomb.io:443"
	}
	ctx := context.Background()
	//creds := credentials.NewClientTLSFromCert(nil, "")
	driver := otlpgrpc.NewDriver(
		otlpgrpc.WithInsecure(),
		otlpgrpc.WithEndpoint(otlpendpoint))
	exporter, err := otlp.NewExporter(ctx, driver)
	if err != nil {
		log.Fatal(err)
	}
	propagator := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
	otel.SetTextMapPropagator(propagator)
	otel.SetTracerProvider(
		trace.NewTracerProvider(
			trace.WithSpanProcessor(trace.NewBatchSpanProcessor(exporter)),
		),
	)
}

 */
