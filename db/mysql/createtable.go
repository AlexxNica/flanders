package mysql

import (
	"fmt"

	"github.com/weave-lab/flanders/log"
)

const (
	tablePrefix = "messages"
	tableSchema = `CREATE TABLE IF NOT EXISTS %s_%s (
		id bigint(20) unsigned NOT NULL AUTO_INCREMENT,
		generated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		date timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
		micro_ts bigint(18) NOT NULL DEFAULT '0',
		method varchar(50) NOT NULL DEFAULT '',
		reply_reason varchar(100) NOT NULL DEFAULT '',
		ruri varchar(200) NOT NULL DEFAULT '',
		ruri_user varchar(100) NOT NULL DEFAULT '',
		ruri_domain varchar(150) NOT NULL DEFAULT '',
		from_user varchar(100) NOT NULL DEFAULT '',
		from_domain varchar(150) NOT NULL DEFAULT '',
		from_tag varchar(128) NOT NULL DEFAULT '',
		to_user varchar(100) NOT NULL DEFAULT '',
		to_domain varchar(150) NOT NULL DEFAULT '',
		to_tag varchar(128) NOT NULL DEFAULT '',
		pid_user varchar(100) NOT NULL DEFAULT '',
		contact_user varchar(120) NOT NULL DEFAULT '',
		auth_user varchar(120) NOT NULL DEFAULT '',
		callid varchar(120) NOT NULL DEFAULT '',
		callid_aleg varchar(120) NOT NULL DEFAULT '',
		via_1 varchar(256) NOT NULL DEFAULT '',
		via_1_branch varchar(80) NOT NULL DEFAULT '',
		cseq varchar(25) NOT NULL DEFAULT '',
		diversion varchar(256) NOT NULL DEFAULT '',
		reason varchar(200) NOT NULL DEFAULT '',
		content_type varchar(256) NOT NULL DEFAULT '',
		auth varchar(256) NOT NULL DEFAULT '',
		user_agent varchar(256) NOT NULL DEFAULT '',
		source_ip varchar(60) NOT NULL DEFAULT '',
		source_port int(10) NOT NULL DEFAULT 0,
		destination_ip varchar(60) NOT NULL DEFAULT '',
		destination_port int(10) NOT NULL DEFAULT 0,
		contact_ip varchar(60) NOT NULL DEFAULT '',
		contact_port int(10) NOT NULL DEFAULT 0,
		originator_ip varchar(60) NOT NULL DEFAULT '',
		originator_port int(10) NOT NULL DEFAULT 0,
		correlation_id varchar(256) NOT NULL DEFAULT '',
		custom_field1 varchar(120) NOT NULL DEFAULT '',
		custom_field2 varchar(120) NOT NULL DEFAULT '',
		custom_field3 varchar(120) NOT NULL DEFAULT '',
		proto int(5) NOT NULL DEFAULT 0,
		family int(1) DEFAULT NULL,
		rtp_stat varchar(256) NOT NULL DEFAULT '',
		type int(2) NOT NULL DEFAULT 0,
		node varchar(125) NOT NULL DEFAULT '',
		msg varchar(3000) NOT NULL DEFAULT '',
		PRIMARY KEY (id,date),
		KEY ruri_user (ruri_user),
		KEY from_user (from_user),
		KEY to_user (to_user),
		KEY pid_user (pid_user),
		KEY auth_user (auth_user),
		KEY callid_aleg (callid_aleg),
		KEY date (date),
		KEY callid (callid),
		KEY method (method)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8 ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8;`
)

// generates the table for the next day
func (m *MySQL) createTable(day string) error {
	log.Info(fmt.Sprintf("generating new table [%s_%s]", tablePrefix, day))
	c := fmt.Sprintf(tableSchema, tablePrefix, day)
	_, err := m.db.Exec(c)
	if err != nil {
		return err
	}
	return nil
}
