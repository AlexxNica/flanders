package db

import (
	"encoding/binary"
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"
)

var Db DbHandler

const (
	DATEFORMAT = "Jan 2, 2006 at 3:04pm (MST)"
)

type Time struct {
	time.Time
}

func (t *Time) SetBSON(raw bson.Raw) error {
	i := int64(binary.LittleEndian.Uint64(raw.Data))
	if i == -62135596800000 {
		t.Time = time.Time{} // In UTC for convenience.
	} else {
		t.Time = time.Unix(i/1e3, i%1e3*1e6)
	}
	return nil
}

func (t Time) MarshalText() ([]byte, error) {
	return []byte(t.Format(DATEFORMAT)), nil
}

func (t Time) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.Format(DATEFORMAT) + `"`), nil
}

type DbObject struct {
	Datetime        Time
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

func (slice DbResult) Len() int {
	return len(slice)
}

func (slice DbResult) Less(i, j int) bool {
	return slice[i].Datetime.Before(slice[j].Datetime.Time)
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
	Sort  []string
	Limit int
}

type DbHandler interface {
	Connect(connectString string) error
	Insert(dbObject *DbObject) error
	Find(filter *Filter, options *Options, result *DbResult) error
	CheckSchema() error // Check to see if the database has been setup or not. Returns nil if all is well
	SetupSchema() error // Sets up the database schema. This will delete all data!!!
}

func RegisterHandler(dbHandler DbHandler) {
	Db = dbHandler
	err := Db.Connect("localhost")
	if err != nil {
		fmt.Println(err)
	}

	_ = Db.SetupSchema()
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
