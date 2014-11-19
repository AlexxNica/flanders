package influx

import (
	"fmt"
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
	a := &db.DbObject{}
	s := reflect.ValueOf(a).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		newInfluxHandler.columns = append(newInfluxHandler.columns, typeOfT.Field(i).Name)
	}
}

func (self *InfluxDb) Connect(connectString string) error {
	var err error
	config := &client.ClientConfig{
		Database: "flanders",
	}
	self.connection, err = client.New(config)
	if err != nil {
		return err
	}

	return nil
}

func (self *InfluxDb) Insert(dbObject *db.DbObject) error {
	newSeries := &client.Series{
		Name:    "message",
		Columns: self.columns,
	}
	var seriesArray []*client.Series
	var values []interface{}

	s := reflect.ValueOf(dbObject).Elem()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		values = append(values, f.Interface())
	}
	newSeries.Points = append(newSeries.Points, values)
	seriesArray = append(seriesArray, newSeries)

	err := self.connection.WriteSeries(seriesArray)
	if err != nil {
		return err
	}
	return nil
}

func (self *InfluxDb) Find(params db.SearchMap, options *db.Options, result *[]db.DbObject) error {

	query := "select * from message"
	count := 0

	for key, val := range params {
		if val != "" {
			if count > 0 {
				query += " AND "
			} else {
				query += " WHERE "
			}
			query += key + "='" + val + "'"
		}
	}

	results, error := self.connection.Query(query)
	if error != nil {
		fmt.Println(error)
	}

	fmt.Printf("%#v\n", results[0])
	return nil
}
