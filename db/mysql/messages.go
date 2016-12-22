package mysql

import (
	"fmt"
	"strings"

	"github.com/weave-lab/flanders/db"
)

func (m *MySQL) prepareInsertQuery() error {

	q := `INSERT INTO messages (
			date, micro_ts,
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
		   ?,?,?,?,?,?,?,?,?,?
		   )`

	i, err := m.db.Prepare(q)
	if err != nil {
		return err
	}

	m.insert = i

	return nil
}

// Insert adds a new record to the messages table in mysql
func (m *MySQL) Insert(d *db.DbObject) error {

	_, err := m.insert.Exec(
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
		d.Type, d.Node, d.Msg,
	)
	if err != nil {
		return err
	}

	return nil
}

const columns = `date, micro_ts,
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
  			type, node, msg`

var (
	filterMap = map[string]string{
		"touser":        "to_user",
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
			return nil, fmt.Errorf("unsupported filter: %s", f)
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

	where := ""
	if len(filters) > 0 {
		where = "WHERE " + strings.Join(filters, " AND ")
	}

	limit := 1000
	if options.Limit > 0 {
		limit = options.Limit
	}

	var order string
	for i, v := range options.Sort {
		dir := " ASC"
		if strings.HasPrefix(v, "-") {
			strings.TrimPrefix(v, "-")
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

	q := fmt.Sprintf(`SELECT %s FROM messages
						%s
						%s
						LIMIT %d`, columns, where, order, limit)

	rows, err := m.db.Query(q, values...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results db.DbResult
	for rows.Next() {

		var d db.DbObject
		err = rows.Scan(
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

		results = append(results, d)
	}

	return results, nil
}
