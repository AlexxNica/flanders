package db

import (
	"fmt"
)

var db DbHandler

type DbObject struct {
	Timestamp       uint32
	TimestampMicro  uint32
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

type SearchMap map[string]interface{}
type OptionsMap struct {
	Sort  []string
	Limit uint
}

type DbHandler interface {
	Connect(connectString string) error
	Insert(dbObject *DbObject) error
	Find(params SearchMap, options OptionsMap, result []DbObject) error
}

func RegisterHandler(dbHandler DbHandler) {
	db = dbHandler
	err := db.Connect("localhost")
	if err != nil {
		fmt.Println(err)
	}
}

func NewDbObject() *DbObject {
	newDbObject := &DbObject{}
	return newDbObject
}

func (d *DbObject) Save() error {
	err := db.Insert(d)
	if err != nil {
		return err
	}
	return nil
}
