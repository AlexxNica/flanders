package db

import (
	"fmt"
	"time"
)

var Db DbHandler
var dbConnectString *string

const (
	DATEFORMAT = "1/2/2006 3:04:05.000pm (MST)"
)

// type Time struct {
// 	time.Time
// }

// func (t *Time) SetBSON(raw bson.Raw) error {
// 	i := int64(binary.LittleEndian.Uint64(raw.Data))
// 	if i == -62135596800000 {
// 		t.Time = time.Time{} // In UTC for convenience.
// 	} else {
// 		t.Time = time.Unix(i/1e3, i%1e3*1e6)
// 	}
// 	return nil
// }

// func (t Time) MarshalText() ([]byte, error) {
// 	return []byte(t.Format(DATEFORMAT)), nil
// }

// func (t Time) MarshalJSON() ([]byte, error) {
// 	return []byte(`"` + t.Format(DATEFORMAT) + `"`), nil
// }

type DbObject struct {
	Datetime        time.Time
	MicroSeconds    int
	Method          string
	ReplyReason     string
	Ruri            string
	RuriUser        string
	RuriDomain      string
	FromUser        string
	FromDomain      string
	FromTag         string
	ToUser          string
	ToDomain        string
	ToTag           string
	PidUser         string
	ContactUser     string
	AuthUser        string
	CallId          string
	CallIdAleg      string
	Via             string
	ViaBranch       string
	Cseq            string
	Diversion       string
	Reason          string
	ContentType     string
	Auth            string
	UserAgent       string
	SourceIp        string
	SourcePort      uint16
	DestinationIp   string
	DestinationPort uint16
	ContactIp       string
	ContactPort     uint16
	OriginatorIp    string
	OriginatorPort  uint16
	Proto           uint
	Family          uint
	RtpStat         string
	Type            uint
	Node            string
	Msg             string
}

type DbResult []DbObject

type SettingObject struct {
	Key string `param:"key"`
	Val string `param:"val"`
}

type SettingResult []SettingObject

// func init() {
// 	dbConnectString = flag.String("db", "localhost", "Database Connection String")
// }

func (slice DbResult) Len() int {
	return len(slice)
}

func (slice DbResult) Less(i, j int) bool {
	if slice[i].Datetime.Equal(slice[j].Datetime) {
		return slice[i].MicroSeconds < slice[j].MicroSeconds
	}

	return slice[i].Datetime.Before(slice[j].Datetime)
}

func (slice DbResult) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

type Filter struct {
	StartDate string
	EndDate   string
	Equals    map[string]interface{}
	Like      map[string]string
}

func NewFilter() Filter {
	filter := Filter{}
	filter.Equals = make(map[string]interface{})
	filter.Like = make(map[string]string)
	return filter
}

type Options struct {
	Sort     []string
	Limit    int
	Distinct string
}

type DbHandler interface {
	Connect(connectString string) error
	Insert(dbObject *DbObject) error
	Find(filter *Filter, options *Options, result *DbResult) error
	CheckSchema() error // Check to see if the database has been setup or not. Returns nil if all is well
	SetupSchema() error // Sets up the database schema. This will delete all data!!!
	GetSettings(settingtype string, result *SettingResult) error
	SetSetting(settingtype string, setting SettingObject) error
	DeleteSetting(settingtype string, key string) error
}

var dbs = make(map[string]DbHandler)

func RegisterHandler(name string, dbHandler DbHandler) {
	dbs[name] = dbHandler
}

func Setup(name string, address string) error {

	var ok bool
	Db, ok = dbs[name]
	if !ok {
		return fmt.Errorf("unknown db: %s", name)
	}

	err := Db.Connect(address)
	if err != nil {
		return fmt.Errorf("unable to connect to db at %s: %s", address, err)
	}

	err = Db.SetupSchema()
	if err != nil {
		return fmt.Errorf("unable to setup db: %s", err)
	}

	return nil
}

func NewDbObject() *DbObject {
	newDbObject := &DbObject{}
	return newDbObject
}

func (d *DbObject) Save() error {
	err := Db.Insert(d)
	if err != nil {
		return err
	}
	return nil
}

// func (d DbObject) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(struct {
// 		DbObject
// 		Datetime string
// 	}{
// 		DbObject: d,
// 		Datetime: d.Datetime.Format(DATEFORMAT),
// 	})
// }
