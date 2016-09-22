package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // Loading mysql driver for this database connection
	"github.com/weave-lab/flanders/db"
)

const (
	dbName         = "sipcapture"
	dataExpiration = 7 // in Days
)

type MysqlDB struct {
	connection *sql.DB
}

func init() {
	newMysqlDB := &MysqlDB{}
	db.RegisterHandler("mysql", newMysqlDB)
}

func (m *MysqlDB) Connect(connectString string) error {
	connection, err := sql.Open("mysql", connectString)
	m.connection = connection
	return err
}
func (m *MysqlDB) Insert(dbObject *db.DbObject) error {
	return nil
}
func (m *MysqlDB) Find(filter *db.Filter, options *db.Options, result *db.DbResult) error {
	return nil
}
func (m *MysqlDB) CheckSchema() error {
	return nil
}
func (m *MysqlDB) SetupSchema() error {
	return nil
}
func (m *MysqlDB) GetSettings(settingtype string, result *db.SettingResult) error {
	return nil
}
func (m *MysqlDB) SetSetting(settingtype string, setting db.SettingObject) error {
	return nil
}
func (m *MysqlDB) DeleteSetting(settingtype string, key string) error {
	return nil
}
