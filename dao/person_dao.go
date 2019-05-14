package dao

import (
	"log"
	"time"

	. "github.com/tadamhicks/rest-api/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// The struct for a person
type PersonDAO struct {
	Server   string
	Port     string
	Database string
	Username string
	Password string
}

var db *mgo.Database

// Collection of struct
const (
	COLLECTION = "person"
)

// Connect is the mgo DB connection method
func (m *PersonDAO) Connect() {

	dialInfo := &mgo.DialInfo{
		Addrs:    []string{m.Server},
		Timeout:  60 * time.Second,
		Database: m.Database,
		Username: m.Username,
		Password: m.Password,
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

// FindAll finds all under /people
func (m *PersonDAO) FindAll() ([]Person, error) {
	var person []Person
	err := db.C(COLLECTION).Find(bson.M{}).All(&person)
	return person, err
}

// FindById returns a single /people by id under /people/{id}
func (m *PersonDAO) FindById(id string) (Person, error) {
	var person Person
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&person)
	return person, err
}

// Insert is the POST endpoint for /people for creating new db entries
func (m *PersonDAO) Insert(person Person) error {
	err := db.C(COLLECTION).Insert(&person)
	return err
}

// Delete is the DELETE endpoint for /people/{id} for removal
func (m *PersonDAO) Delete(id string) error {
	err := db.C(COLLECTION).Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	return err
}

// Update is the PUT endpoint for /people/{id}
func (m *PersonDAO) Update(id string, person Person) error {
	err := db.C(COLLECTION).Update(bson.M{"_id": bson.ObjectIdHex(id)}, &person)
	return err
}
