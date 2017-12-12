package mysql

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/weave-lab/flanders/db"
	"github.com/weave-lab/flanders/log"
)

// Insert adds a new record to the messages table in mysql
func (m *MySQL) Insert(d db.DbObject) error {
	if *batchInsert {
		m.batch.Lock()
		m.batch.rows = append(m.batch.rows, d)
		m.batch.Unlock()
		if len(m.batch.rows) < m.batch.maxRows {
			return nil
		}
		return m.processBatch(m.batch.rows)
	}

	//gzip full packet
	var gzMsg bytes.Buffer
	w := gzip.NewWriter(&gzMsg)
	w.Write([]byte(d.Msg))
	w.Close()

	_, err := m.insert[time.Now().Format("01_02_2006")].Exec(
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

	var tables []string
	var err error
	var startDate time.Time
	var endDate time.Time

	if filter.StartDate != "" {
		log.Debug(fmt.Sprintf("filtering by start date [%s]", filter.StartDate))
		filters = append(filters, "date >= ?")
		values = append(values, filter.StartDate)
		startDate, err = time.Parse(time.RFC3339, filter.StartDate)
		if err != nil {
			log.Info(fmt.Sprintf("could not parse start date [%s]", err.Error()))
		}
	}

	if filter.EndDate != "" {
		log.Debug(fmt.Sprintf("filtering by end date [%s]", filter.EndDate))
		filters = append(filters, "date < ?")
		values = append(values, filter.EndDate)
		endDate, err = time.Parse(time.RFC3339, filter.EndDate)
		if err != nil {
			log.Info(fmt.Sprintf("could not parse end date [%s]", err.Error()))
		}
	}


	if startDate.IsZero() && endDate.IsZero() { // only use today
		tables = append(tables, fmt.Sprintf("%s_%s", tablePrefix, time.Now().Format("01_02_2006")))
	} else if !startDate.IsZero() && !endDate.IsZero() && endDate.After(startDate) {
		for d := startDate; !d.Equal(endDate.AddDate(0, 0, 1)); d = d.AddDate(0, 0, 1) {
			date := d.Format("01_02_2006")
			log.Debug(fmt.Sprintf("adding date [%s]", date))
			tables = append(tables, fmt.Sprintf("%s_%s", tablePrefix, date))
		}
	} else if !startDate.IsZero() {
		date := startDate.Format("01_02_2006")
		log.Debug(fmt.Sprintf("adding date [%s]", date))
		tables = append(tables, fmt.Sprintf("%s_%s", tablePrefix, date))
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

	// aggregate the data from all tables
	var results db.DbResult
	for _, table := range tables {
		q := fmt.Sprintf(`SELECT %s FROM %s %s %s LIMIT %d`, strings.Join(columns, ","), table, where, order, limit)
		//if options.UniqueCallID {
			// Wrap query in another query to group by call ID instead of returning
			// every packet. We just want individual calls
			// THIS CAUSES "[Error 1055: Expression #1 of SELECT list is not in GROUP BY clause and contains nonaggregated column 'x.generated_at' which is not functionally dependent on columns in GROUP BY clause; this is incompatible with sql_mode=only_full_group_by]"
		//	q = fmt.Sprintf(`SELECT %s FROM (%s) x GROUP BY x.callid`, strings.Join(columns, ","), q)
		//}
		log.Debug(fmt.Sprintf("executing query [%s]", q))
		rows, err := m.db.Query(q, values...)
		if err != nil {
			log.Crit(fmt.Sprintf("could not execute query [%s]", err.Error()))
			continue
		}
		defer rows.Close()

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
				//log.Crit(err.Error())
			} else {
				uzipMsg, err := ioutil.ReadAll(r)
				r.Close()
				if err != nil {
					//fmt.Println(err)
					//return nil, err
				} else {
					d.Msg = string(uzipMsg)
				}
			}

			results = append(results, d)
		}

	}

	return results, nil
}
