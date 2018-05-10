package dao

import (
	"log"
	. "github.com/tadamhicks/rest-api/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type PersonDAO struct {
	Server		string
	Database	string
	Username	string
	Password	string
}

var db *mgo.Database

const (
	COLLECTION = "person"
)

func (m *PersonDAO) Connect() {

	dialInfo := &mgo.DialInfo {
		Addrs: []string{m.Server},
		Timeout: 60 * time.Second,
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


func (m *PersonDAO) FindAll() ([]Person, error) {
	var person []Person
	err := db.C(COLLECTION).Find(bson.M{}).All(&person)
	return person, err
}


func (m *PersonDAO) FindById(id string) (Person, error) {
	var person Person
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&person)
	return person, err
}


func (m *PersonDAO) Insert(person Person) error {
	err := db.C(COLLECTION).Insert(&person)
	return err
}


func (m *PersonDAO) Delete(id string) error {
	err := db.C(COLLECTION).Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	return err
}


func (m *PersonDAO) Update(id string, person Person) error {
	err := db.C(COLLECTION).Update(bson.M{"_id": bson.ObjectIdHex(id)}, &person)
	return err
}
