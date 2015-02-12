package influx

import (
	"fmt"
	"reflect"
	s "strings"

	"github.com/influxdb/influxdb/client"
	"github.com/weave-lab/flanders/db"
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

func (self *InfluxDb) Find(filter *db.Filter, options *db.Options, result *[]db.DbObject) error {

	query := "select * from message"
	conditions := make([]string, 0)

	if filter.StartDate != "" {
		conditions = append(conditions, "time > '"+filter.StartDate+"'")
	}
	if filter.EndDate != "" {
		conditions = append(conditions, "time < '"+filter.EndDate+"'")
	}

	for key, val := range filter.Equals {
		if val != "" {
			query += key + " = '" + val + "'"
		}
	}

	for key, val := range filter.Like {
		if val != "" {
			query += key + " LIKE '" + val + "'"
		}
	}

	if len(conditions) > 0 {
		query += " WHERE "
		query += s.Join(conditions, " AND ")
	}

	results, error := self.connection.Query(query)
	if error != nil {
		fmt.Println(error)
	}

	for key, val := range results {
		newDbObject := &db.DbObject{}
	}

	fmt.Printf("%#v\n", len(results))
	return nil
}
