package mongo

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"lab.getweave.com/weave/flanders/db"
	"time"
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

func (m *MongoDb) Find(filter *db.Filter, options *db.Options, result *[]db.DbObject) error {
	collection := m.connection.DB(DB_NAME).C("message")

	conditions := bson.M{}
	var err error
	var startDate time.Time
	var endDate time.Time

	if filter.StartDate != "" {
		fmt.Print("Start date found... " + filter.StartDate)
		startDate, err = time.Parse(time.RFC3339, filter.StartDate)
		if err != nil {
			return errors.New("Could not parse `Start Date` from filters")
		}
		conditions["datetime"] = bson.M{"$gte": startDate}
	}
	if filter.EndDate != "" {
		fmt.Print("End date found... " + filter.EndDate)
		endDate, err = time.Parse(time.RFC3339, filter.EndDate)
		if err != nil {
			return errors.New("Could not parse `End Date` from filters")
		}
		conditions["datetime"] = bson.M{"$lt": endDate}
	}
	for key, val := range filter.Equals {
		conditions[key] = val
	}

	for key, val := range filter.Like {
		conditions[key] = "/" + val + "/"
	}

	query := collection.Find(conditions)

	//sort := options.Sort

	// if sort != nil {
	// 	query = query.Sort(...sort)
	// } else {
	// 	query = query.Sort("Timestamp")
	// }

	query.All(result)

	return nil
}

func (self *MongoDb) CheckSchema() error {
	return nil
}

func (self *MongoDb) SetupSchema() error {
	return nil
}
