package mysql

import (
	"database/sql"
	"fmt"

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

func (m *MySQL) CheckSchema() error {
	rows, err := m.db.Query(`SELECT count(*) FROM messages`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		return nil
	}
	return fmt.Errorf("schema not found")
}

func (m *MySQL) SetupSchema() error {
	return fmt.Errorf("setup schema is not implemented")
}
