package influx

import (
	//"fmt"
	"github.com/influxdb/influxdb/client"
	"lab.getweave.com/weave/flanders/db"
	"reflect"
)

const (
	DB_NAME = "flanders"
)

type InfluxDb struct {
	connection *client.Client
	columns    []string
}

func init() {
	newInfluxHandler := &InfluxDb{}
	db.RegisterHandler(newInfluxHandler)

	// Get DbObject struct field names for columns object
	a := &b.DbObject{}
	s := reflect.ValueOf(a).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		typeOfT.Field(i).Name
	}
}

func (i *InfluxDb) Connect(connectString string) error {
	var err error
	config := &client.ClientConfig{
		Database: "flanders",
	}
	i.connection, err = client.New(config)
	if err != nil {
		return err
	}

	return nil
}

func (i *InfluxDb) Insert(dbObject *db.DbObject) error {

	return nil
}

func (i *InfluxDb) Find(params db.SearchMap, options *db.Options, result *[]db.DbObject) error {
	return nil
}
