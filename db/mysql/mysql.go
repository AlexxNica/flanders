package mysql

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" // Loading mysql driver for this database connection
	"github.com/weave-lab/flanders/db"
	"github.com/weave-lab/flanders/log"
)

var (
	maxConnections  *int
	batchInsert     *bool
	batchAmount     *int
	batchFrequency  *int
	cleanOldExpired *int
	cleanUpOnStart  *bool

	insertStatement = `INSERT INTO %s_%s (
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
)

func init() {
	// Add mysql specific flag settings
	maxConnections = flag.Int("mysql-max-connections", 30, "Max connections for mysql")
	cleanOldExpired = flag.Int("sip-message-expire", 5, "clean up sip messages older than X days")
	cleanUpOnStart = flag.Bool("clean-up-on-start", false, "clean up old sip messages on start")
	batchInsert = flag.Bool("mysql-batch", true, "Use batch inserting for high traffic systems")
	batchAmount = flag.Int("mysql-batch-count", 100, "Amount of messages to batch at one time")
	batchFrequency = flag.Int("mysql-batch-frequency", 5, "send batch every X seconds if max it not hit")
}

type MySQL struct {
	db *sql.DB

	insert map[string]*sql.Stmt
	batch  *batch
}

func init() {
	m := MySQL{}
	m.insert = make(map[string]*sql.Stmt)
	b := &batch{}
	b.maxRows = *batchAmount
	m.batch = b

	db.RegisterHandler("mysql", &m)
}

// Connect connects to the database... go figure
func (m *MySQL) Connect(connectString string) error {
	connection, err := sql.Open("mysql", connectString)
	connection.SetMaxOpenConns(*maxConnections)
	if err != nil {
		return err
	}

	m.db = connection

	err = m.SetupSchema()
	if err != nil {
		m.db.Close()
		return err
	}

	err = m.prepareInserts()
	if err != nil {
		m.db.Close()
		return err
	}
	log.Debug("connected to mysql")

	if *batchInsert {
		log.Info(fmt.Sprintf("batch enabled with max rows [%d]", m.batch.maxRows))
		go m.runBatch()
	}

	// cron job to clean up old sip messages
	m.startCron()

	if *cleanUpOnStart {
		go func() {
			err := m.cleanup()
			if err != nil {
				log.Crit(fmt.Sprintf("could not clean up messages on start [%s]", err.Error()))
			}
		}()
	}

	return nil
}

func (m *MySQL) runBatch() error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sendFrequency := time.Duration(*batchFrequency) * time.Second
	t := time.NewTicker(sendFrequency)
	for {
		select {
		case <-sigChan:
			log.Info("received shutdown message")
			if len(m.batch.rows) < m.batch.maxRows {
				return nil
			}
			log.Info(fmt.Sprintf("sending final batch of rows [%d]", len(m.batch.rows)))
			return m.processBatchOnTimer(m.batch.rows)
		case <-t.C:
			// call send if it has been > sendFrequency since last send
			if m.batch.lastSent.Add(sendFrequency).Before(time.Now()) {
				err := m.processBatchOnTimer(m.batch.rows)
				if err != nil {
					log.Crit(fmt.Sprintf("could not process batch [%s]", err.Error()))
				}
			}
		}
	}
	return nil
}

func (m *MySQL) prepareInsert(day string) error {
	log.Info(fmt.Sprintf("preparing insert statement for table [%s_%s]", tablePrefix, day))
	q := fmt.Sprintf(insertStatement, tablePrefix, day)
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

	m.insert[day] = i
	return nil
}

func (m *MySQL) prepareInserts() error {
	day := fmt.Sprintf(time.Now().Format("01_02_2006"))
	err := m.prepareInsert(day)
	if err != nil {
		return err
	}
	day = fmt.Sprintf(time.Now().Add(time.Hour * 24).Format("01_02_2006")) // prepare statement for next day too
	err = m.prepareInsert(day)
	if err != nil {
		return err
	}
	day = fmt.Sprintf(time.Now().Add(-time.Hour * 48).Format("01_02_2006")) // remove old insert statements
	log.Info(fmt.Sprintf("deleting old insert statement [%s_%s]", tablePrefix, day))
	delete(m.insert, day)
	return nil
}

//CheckSchema checks to make sure that the database schema will work with this version of Flanders
func (m *MySQL) CheckSchema() error {
	today := fmt.Sprintf(time.Now().Format("01_02_2006"))
	table := fmt.Sprintf("%s_%s", tablePrefix, today)

	q := fmt.Sprintf("SELECT date FROM %s LIMIT 10;", table)
	rows, err := m.db.Query(q)
	if err != nil {
		return err
	}
	defer rows.Close()

	type schemaTester struct {
		Datetime time.Time
	}

	for rows.Next() {
		var st schemaTester
		err = rows.Scan(
			&st.Datetime,
		)
		if err != nil {
			switch {
			case strings.Contains(err.Error(), "Scan error on column index 0"):
				return fmt.Errorf("schema error parsing Date.  Did you include the DSN parameter parseTime=true on your connection string?")
			}
			return fmt.Errorf("schema error %s", err)
		}
	}
	return nil
}

// setup schemas
func (m *MySQL) SetupSchema() error {
	day := fmt.Sprintf(time.Now().Format("01_02_2006"))
	err := m.createTable(day)
	if err != nil {
		return err
	}
	day = fmt.Sprintf(time.Now().Add(time.Hour * 24).Format("01_02_2006"))
	err = m.createTable(day)
	if err != nil {
		return err
	}
	return nil
}
