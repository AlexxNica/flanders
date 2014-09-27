package mongo

import (
	//"fmt"
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
	"lab.getweave.com/weave/flanders/db"
)

const (
	DB_NAME = "flanders"
)

type MongoDb struct {
	connection *mgo.Session
}

func init() {
	newMongoHandler := &MongoDb{}
	db.RegisterHandler(newMongoHandler)
}

func (m *MongoDb) Connect(connectString string) error {
	var err error
	m.connection, err = mgo.Dial(connectString)
	if err != nil {
		return err
	}

	// Optional. Switch the connection to a monotonic behavior.
	m.connection.SetMode(mgo.Monotonic, true)
	return nil
}

func (m *MongoDb) Insert(dbObject *db.DbObject) error {
	collection := m.connection.DB(DB_NAME).C("message")
	err := collection.Insert(dbObject)
	return err
}

func (m *MongoDb) Find(params db.SearchMap, options db.OptionsMap, result []db.DbObject) error {
	collection := m.connection.DB(DB_NAME).C("message")
	query := collection.Find(params)

	var sort []string
	var okSort bool

	sort, okSort = options["sort"]
	if okSort {
		query = query.Sort(sort...)
	} else {
		query = query.Sort("Timestamp")
	}
}
