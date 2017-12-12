package mysql

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"sync"
	"time"

	"github.com/weave-lab/flanders/db"
	"github.com/weave-lab/flanders/log"
)

type batch struct {
	sync.Mutex
	maxRows  int
	lastSent time.Time
	rows     []db.DbObject
}

func (m *MySQL) processBatch(rows []db.DbObject) error {
	m.batch.Lock()
	defer m.batch.Unlock()
	if len(rows) == 0 {
		log.Debug("ignoring batch with 0 rows")
		return nil
	}
	m.batch.rows = []db.DbObject{} // reset batch
	go func() {
		err := m.insertBatch(rows)
		if err != nil {
			log.Crit(fmt.Sprintf("could not process batch [%s]", err.Error()))
		}
	}()
	m.batch.lastSent = time.Now()
	return nil
}

// insertBatch inserts a group of sip messages
func (m *MySQL) insertBatch(rows []db.DbObject) error {
	var hugeInsertSlice []interface{}

	for _, d := range rows {
		//gzip full packet
		var gzMsg bytes.Buffer
		w := gzip.NewWriter(&gzMsg)
		w.Write([]byte(d.Msg))
		w.Close()
		tempSlice := []interface{}{
			d.GeneratedAt, d.Datetime, d.MicroSeconds,
			d.Method, d.ReplyReason, d.Ruri,
			d.RuriUser, d.RuriDomain,
			d.FromUser, d.FromDomain, d.FromTag,
			d.ToUser, d.ToDomain, d.ToTag,
			d.PidUser, d.ContactUser, d.AuthUser,
			d.CallId, d.CallIdAleg,
			d.Via, d.ViaBranch,
			d.Cseq, d.Diversion,
			d.Reason, d.ContentType,
			d.Auth, d.UserAgent,
			d.SourceIp, d.SourcePort,
			d.DestinationIp, d.DestinationPort,
			d.ContactIp, d.ContactPort,
			d.OriginatorIp, d.OriginatorPort,
			d.Proto, d.Family, d.RtpStat,
			d.Type, d.Node, gzMsg.String(),
		}

		hugeInsertSlice = append(hugeInsertSlice, tempSlice...)

	}
	_, err := m.insert[time.Now().Format("01_02_2006")].Exec(hugeInsertSlice...)

	return err
}
