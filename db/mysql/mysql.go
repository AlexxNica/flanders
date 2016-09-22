package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // Loading mysql driver for this database connection
	"github.com/vattle/sqlboiler/queries/qm"

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
	dbObjectToMessage(dbObject, message)
	return message.Insert(m.connection)
}

func (m *MysqlDB) Find(filter *db.Filter, options *db.Options, result *db.DbResult) error {
	// var startDate time.Time
	// var endDate time.Time

	var conditions []qm.QueryMod

	// if filter.StartDate != "" {
	// 	log.Debug("Start date found... " + filter.StartDate)
	// 	startDate, err := time.Parse(time.RFC3339, filter.StartDate)
	// 	if err != nil {
	// 		return errors.New("Could not parse `Start Date` from filters")
	// 	}
	// 	conditions = append(conditions, qm.Where("date>=?", startDate))
	// }
	// if filter.EndDate != "" {
	// 	log.Debug("End date found... " + filter.EndDate)
	// 	endDate, err := time.Parse(time.RFC3339, filter.EndDate)
	// 	if err != nil {
	// 		return errors.New("Could not parse `End Date` from filters")
	// 	}
	// 	conditions = append(conditions, qm.Where("date<=?", endDate))
	// }
	conditions = append(conditions, qm.OrderBy("date DESC"), qm.Limit(100))
	messages, err := models.Messages(m.connection, conditions...).All()
	if err != nil {
		return err
	}

	messagesToDbResult(messages, result)

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

func messageToDbObject(message *models.Message, dbObject *db.DbObject) {
	dbObject.Datetime = message.Date
	dbObject.MicroSeconds = int(message.MicroTS)
	dbObject.Method = message.Method
	dbObject.ReplyReason = message.ReplyReason
	dbObject.SourceIp = message.SourceIp
	dbObject.SourcePort = uint16(message.SourcePort)
	dbObject.DestinationIp = message.DestinationIp
	dbObject.DestinationPort = uint16(message.DestinationPort)
	dbObject.CallId = message.Callid
	dbObject.FromUser = message.FromUser
	dbObject.FromDomain = message.FromDomain
	dbObject.FromTag = message.FromTag
	dbObject.ToUser = message.ToUser
	dbObject.ToDomain = message.ToDomain
	dbObject.ToTag = message.ToTag
	dbObject.UserAgent = message.UserAgent
	dbObject.Cseq = message.Cseq
	dbObject.Msg = message.MSG
}

func dbObjectToMessage(dbObject *db.DbObject, message *models.Message) {
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
	message.MSG = dbObject.Msg
}

func messagesToDbResult(messages models.MessageSlice, result *db.DbResult) {
	for _, m := range messages {
		d := &db.DbObject{}
		messageToDbObject(m, d)
		result.Append(d)
	}
}
