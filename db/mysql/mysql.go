package mysql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" // Loading mysql driver for this database connection
	"github.com/weave-lab/flanders/db"
)

type MySQL struct {
	db *sql.DB

	insert *sql.Stmt
}

func init() {
	m := MySQL{}

	db.RegisterHandler("mysql", &m)
}

func (m *MySQL) Connect(connectString string) error {
	connection, err := sql.Open("mysql", connectString)
	connection.SetMaxOpenConns(10)
	if err != nil {
		return err
	}

	m.db = connection

	err = m.prepareInsertQuery()
	if err != nil {
		m.db.Close()
		return err
	}

	return err
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
