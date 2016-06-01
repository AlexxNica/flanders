package mongo

import (
	"errors"
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"lab.getweave.com/weave/flanders/db"
)

const (
	DB_NAME         = "flanders"
	DATA_EXPIRATION = 7 // in Days
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

func (m *MongoDb) Find(filter *db.Filter, options *db.Options, result *db.DbResult) error {
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
		conditions[key] = bson.M{"$regex": bson.RegEx{`\` + val + `\`, ""}}
	}

	query := collection.Find(conditions)

	if options.Limit != 0 {
		query = query.Limit(options.Limit)
	}

	if len(options.Sort) != 0 {
		query = query.Sort(options.Sort...)
	} else {
		query = query.Sort("-datetime")
	}

	//sort := options.Sort

	// if sort != nil {
	// 	query = query.Sort(...sort)
	// } else {
	// 	query = query.Sort("Timestamp")
	// }

	query.All(result)

	return nil
}

func (m *MongoDb) GetSettings(settingtype string, result *db.SettingResult) error {
	collection := m.connection.DB(DB_NAME).C(settingtype)

	query := collection.Find(bson.M{})
	query.All(result)

	return nil
}

func (m *MongoDb) SetSetting(settingtype string, setting db.SettingObject) error {
	collection := m.connection.DB(DB_NAME).C(settingtype)
	_, err := collection.Upsert(bson.M{"key": setting.Key}, setting)

	return err
}

func (m *MongoDb) DeleteSetting(settingtype string, key string) error {
	collection := m.connection.DB(DB_NAME).C(settingtype)
	err := collection.Remove(bson.M{"key": key})

	return err
}

func (m *MongoDb) CheckSchema() error {
	return nil
}

func (m *MongoDb) SetupSchema() error {
	collection := m.connection.DB(DB_NAME).C("message")
	var err error

	callidIndex := mgo.Index{
		Key:        []string{"callid"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     false,
	}
	err = collection.EnsureIndex(callidIndex)
	if err != nil {
		return err
	}

	callidalegIndex := mgo.Index{
		Key:        []string{"callidaleg"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     false,
	}
	err = collection.EnsureIndex(callidalegIndex)
	if err != nil {
		return err
	}

	touserIndex := mgo.Index{
		Key:        []string{"touser"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     false,
	}
	err = collection.EnsureIndex(touserIndex)
	if err != nil {
		return err
	}

	fromuserIndex := mgo.Index{
		Key:        []string{"fromuser"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     false,
	}
	err = collection.EnsureIndex(fromuserIndex)
	if err != nil {
		return err
	}

	fromdomainIndex := mgo.Index{
		Key:        []string{"sourceip"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     false,
	}
	err = collection.EnsureIndex(fromdomainIndex)
	if err != nil {
		return err
	}

	todomainIndex := mgo.Index{
		Key:        []string{"destinationip"},
		Unique:     false,
		DropDups:   false,
		Background: true,
		Sparse:     false,
	}
	err = collection.EnsureIndex(todomainIndex)
	if err != nil {
		return err
	}

	datetimeIndex := mgo.Index{
		Key:         []string{"datetime"},
		Unique:      false,
		DropDups:    false,
		Background:  true,
		Sparse:      false,
		ExpireAfter: time.Duration(DATA_EXPIRATION*24) * time.Hour,
	}
	err = collection.EnsureIndex(datetimeIndex)
	if err != nil {
		return err
	}
	return nil
}
