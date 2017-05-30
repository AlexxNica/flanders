package mysql

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/weave-lab/flanders/db"
)

var batch []*db.DbObject
var lock sync.Mutex

func (m *MySQL) prepareInsertQuery() error {

	q := `INSERT INTO messages (
			generated_at, date, micro_ts,
			method, reply_reason, ruri,
  			ruri_user, ruri_domain,
  			from_user, from_domain, from_tag,
  			to_user, to_domain, to_tag,
  			pid_user, contact_user, auth_user,
  			callid, callid_aleg,
  			via_1, via_1_branch,
  			cseq, diversion,
  			reason, content_type,
  			auth, user_agent,
			source_ip, source_port,
  			destination_ip, destination_port,
  			contact_ip, contact_port,
  			originator_ip, originator_port,
  			proto, family, rtp_stat,
  			type, node, msg
		)
	VALUES(?,?,?,?,?,?,?,?,?,?,
		   ?,?,?,?,?,?,?,?,?,?,
		   ?,?,?,?,?,?,?,?,?,?,
		   ?,?,?,?,?,?,?,?,?,?,
		   ?
		   )`

	if *batchInsert {
		for i := 1; i < *batchAmount; i++ {
			q += `,(?,?,?,?,?,?,?,?,?,?,
    		   ?,?,?,?,?,?,?,?,?,?,
    		   ?,?,?,?,?,?,?,?,?,?,
    		   ?,?,?,?,?,?,?,?,?,?,
    		   ?
    		   )`
		}
	}

	i, err := m.db.Prepare(q)
	if err != nil {
		return err
	}

	m.insert = i

	return nil
}

// Insert adds a new record to the messages table in mysql
func (m *MySQL) Insert(d *db.DbObject) error {

	if *batchInsert {
		lock.Lock()
		batch = append(batch, d)
		if len(batch) == *batchAmount {
			batchForInsert := batch
			batch = nil
			lock.Unlock()
			return m.InsertBatch(batchForInsert)
		}
		lock.Unlock()
		return nil
	}

	//gzip full packet
	var gzMsg bytes.Buffer
	w := gzip.NewWriter(&gzMsg)
	w.Write([]byte(d.Msg))
	w.Close()

	_, err := m.insert.Exec(
		d.GeneratedAt,
		d.Datetime, d.MicroSeconds,
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
	)
	if err != nil {
		return err
	}
	return nil

}

// InsertBatch inserts a group of sip messages
func (m *MySQL) InsertBatch(b []*db.DbObject) error {
	var hugeInsertSlice []interface{}

	for _, d := range b {
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

	_, err := m.insert.Exec(hugeInsertSlice...)

	return err
}

var (
	columns = []string{
		`generated_at`, `date`, `micro_ts`, `method`, `reply_reason`, `ruri`, `ruri_user`, `ruri_domain`,
		`from_user`, `from_domain`, `from_tag`, `to_user`, `to_domain`, `to_tag`, `pid_user`,
		`contact_user`, `auth_user`, `callid`, `callid_aleg`, `via_1`, `via_1_branch`, `cseq`,
		`diversion`, `reason`, `content_type`, `auth`, `user_agent`, `source_ip`, `source_port`,
		`destination_ip`, `destination_port`, `contact_ip`, `contact_port`, `originator_ip`,
		`originator_port`, `proto`, `family`, `rtp_stat`, `type`, `node`, `msg`}

	filterMap = map[string]string{
		"touser":        "to_user",
		"callid":        "callid",
		"callidaleg":    "callid_aleg",
		"fromuser":      "from_user",
		"sourceip":      "source_ip",
		"destinationip": "destination_ip",
	}

	orderMap = map[string]string{
		"datetime":     "date",
		"microseconds": "micro_ts",
	}
)

// Find returns all messages that match the filter parameters
// Options.Distinct is not supported
func (m *MySQL) Find(filter *db.Filter, options *db.Options) (db.DbResult, error) {
	var filters []string
	var orFilters []string
	var values []interface{}

	if filter.StartDate != "" {
		filters = append(filters, "date >= ?")
		values = append(values, filter.StartDate)
	}

	if filter.EndDate != "" {
		filters = append(filters, "date < ?")
		values = append(values, filter.EndDate)
	}

	for k, v := range filter.Equals {
		f, ok := filterMap[k]
		if !ok {
			return nil, fmt.Errorf("unsupported filter: %s [%+v]", k, filter.Equals)
		}

		filters = append(filters, f+" = ?")
		values = append(values, v)
	}

	for k, v := range filter.Like {
		f, ok := filterMap[k]
		if !ok {
			return nil, fmt.Errorf("unsupported filter: %s", f)
		}

		filters = append(filters, f+" LIKE ?")
		values = append(values, v)
	}

	for k, v := range filter.Or {
		f, ok := filterMap[k]
		if !ok {
			return nil, fmt.Errorf("unsupported filter: %s [%+v]", k, filter.Equals)
		}

		orFilters = append(orFilters, f+" = ?")
		values = append(values, v)
	}

	// This is limited but works for existing queries
	// this will NOT work for queries like:
	// WHERE (f1='a' OR f2='b') AND (f3='c' OR f4='d')
	where := ""
	if len(filters) > 0 || len(orFilters) > 0 {
		where = "WHERE "
		if len(filters) > 0 {
			where = where + strings.Join(filters, " AND ")
		}
		if len(orFilters) > 0 {
			if len(filters) > 0 {
				where = where + " AND "
			}
			where = where + " (" + strings.Join(orFilters, " OR ") + ")"
		}
	}

	limit := 1000
	if options.Limit > 0 {
		limit = options.Limit
	}

	var order string
	for i, v := range options.Sort {
		dir := " ASC"
		if strings.HasPrefix(v, "-") {
			v = strings.TrimPrefix(v, "-")
			dir = " DESC"
		}

		s, ok := orderMap[v]
		if !ok {
			return nil, fmt.Errorf("unsupported sort: %s", v)
		}

		comma := ", "
		if i == 0 {
			comma = "ORDER BY "
		}
		order += comma + s + dir
	}
	q := fmt.Sprintf(`SELECT %s
					  FROM messages
					  %s
					  %s
					  LIMIT %d`, strings.Join(columns, ","), where, order, limit)

	if options.UniqueCallID {
		// Wrap query in another query to group by call ID instead of returning
		// every packet. We just want individual calls
		q = fmt.Sprintf(`SELECT * FROM (%s) x GROUP BY x.callid`, q)
	}

	fmt.Println(q)

	rows, err := m.db.Query(q, values...)
	if err != nil {
		fmt.Println(q)
		return nil, err
	}

	defer rows.Close()

	var results db.DbResult
	for rows.Next() {
		var d db.DbObject
		err = rows.Scan(
			&d.GeneratedAt,
			&d.Datetime, &d.MicroSeconds,
			&d.Method, &d.ReplyReason, &d.Ruri,
			&d.RuriUser, &d.RuriDomain,
			&d.FromUser, &d.FromDomain, &d.FromTag,
			&d.ToUser, &d.ToDomain, &d.ToTag,
			&d.PidUser, &d.ContactUser, &d.AuthUser,
			&d.CallId, &d.CallIdAleg,
			&d.Via, &d.ViaBranch,
			&d.Cseq, &d.Diversion,
			&d.Reason, &d.ContentType,
			&d.Auth, &d.UserAgent,
			&d.SourceIp, &d.SourcePort,
			&d.DestinationIp, &d.DestinationPort,
			&d.ContactIp, &d.ContactPort,
			&d.OriginatorIp, &d.OriginatorPort,
			&d.Proto, &d.Family, &d.RtpStat,
			&d.Type, &d.Node, &d.Msg,
		)
		if err != nil {
			return nil, err
		}

		r, err := gzip.NewReader(strings.NewReader(d.Msg))
		if err != nil {
			return nil, err
		}
		uzipMsg, err := ioutil.ReadAll(r)
		r.Close()
		if err != nil {
			return nil, err
		}
		d.Msg = string(uzipMsg)

		results = append(results, d)
	}
	return results, nil
}
