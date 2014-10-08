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
	a := &db.DbObject{}
	s := reflect.ValueOf(a).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		newInfluxHandler.columns = append(newInfluxHandler.columns, typeOfT.Field(i).Name)
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

func (influx *InfluxDb) Insert(dbObject *db.DbObject) error {
	newSeries := &client.Series{
		Name:    "message",
		Columns: influx.columns,
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

	err := influx.connection.WriteSeries(seriesArray)
	if err != nil {
		return err
	}
	return nil
}

func (i *InfluxDb) Find(params db.SearchMap, options *db.Options, result *[]db.DbObject) error {
	return nil
}
