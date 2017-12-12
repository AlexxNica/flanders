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
)

func init() {
	// Add mysql specific flag settings
	maxConnections = flag.Int("mysql-max-connections", 30, "Max connections for mysql")
	cleanOldExpired = flag.Int("sip-message-expire", 2, "clean up sip messages older than X days")
	cleanUpOnStart = flag.Bool("clean-up-on-start", false, "clean up old sip messages on start")
	batchInsert = flag.Bool("mysql-batch", true, "Use batch inserting for high traffic systems")
	batchAmount = flag.Int("mysql-batch-count", 100, "Amount of messages to batch at one time")
	batchFrequency = flag.Int("mysql-batch-frequency", 5, "send batch every X seconds if max it not hit")
}

type MySQL struct {
	db *sql.DB

	insert *sql.Stmt
	batch  *batch
}

func init() {
	m := MySQL{}
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

	err = m.prepareInsertQuery()
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
			return m.processBatch(m.batch.rows)
		case <-t.C:
			log.Info("processing batch on timer")
			// call send if it has been > sendFrequency since last send
			if m.batch.lastSent.Add(sendFrequency).Before(time.Now()) {
				err := m.processBatch(m.batch.rows)
				if err != nil {
					log.Crit(fmt.Sprintf("could not process batch [%s]", err.Error()))
				}
			}
		}
	}
	return nil
}

//CheckSchema checks to make sure that the database schema will work with this version of Flanders
func (m *MySQL) CheckSchema() error {
	rows, err := m.db.Query(`SELECT date FROM messages LIMIT 10;`)
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

func (m *MySQL) SetupSchema() error {
	return fmt.Errorf("setup schema is not implemented")
}
