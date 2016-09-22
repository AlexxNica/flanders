package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // Loading mysql driver for this database connection
	"github.com/weave-lab/flanders/db"
	"github.com/weave-lab/flanders/db/mysql/models"
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

// Insert adds a new record to the messages table in mysql
func (m *MysqlDB) Insert(dbObject *db.DbObject) error {
	message := &models.Message{}
	message.Date = dbObject.Datetime
	message.MicroTS = int64(dbObject.MicroSeconds)
	message.Method = dbObject.Method
	message.ReplyReason = dbObject.ReplyReason
	message.SourceIp = dbObject.SourceIp
	message.SourcePort = int(dbObject.SourcePort)
	message.DestinationIp = dbObject.DestinationIp
	message.DestinationPort = int(dbObject.DestinationPort)
	message.Callid = dbObject.CallId
	message.FromUser = dbObject.FromUser
	message.FromDomain = dbObject.FromDomain
	message.FromTag = dbObject.FromTag
	message.ToUser = dbObject.ToUser
	message.ToDomain = dbObject.ToDomain
	message.ToTag = dbObject.ToTag
	message.UserAgent = dbObject.UserAgent
	message.Cseq = dbObject.Cseq

	return message.Insert(m.connection)
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
