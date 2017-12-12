package mysql

import (
	"fmt"
	"time"

	"github.com/robfig/cron"
	"github.com/weave-lab/flanders/log"
)

var (
	retention time.Duration
)

// use cron to generate new tables each day
func (m *MySQL) startCron() {
	retention = time.Hour * 24 * time.Duration(*cleanOldExpired)
	c := cron.New()
	c.AddFunc("0 0 12 * * *", func() {
		nextDay := fmt.Sprintf(time.Now().Add(time.Hour * 24 + 1).Format("01_02_2006"))
		err := m.createTable(nextDay)
		if err != nil {
			log.Crit(fmt.Sprintf("could not create new table [%s]", err.Error()))
		}
		err = m.prepareInserts()
		if err != nil {
			log.Crit(fmt.Sprintf("could not prepare insert statements [%s]", err.Error()))
		}
	})
	c.AddFunc("0 0 7 * * *", func() {
		err := m.cleanup()
		if err != nil {
			log.Crit(fmt.Sprintf("could not clear old messages [%s]", err.Error()))
		}
	})
	c.Start()
}

//cleanup removes old sip messages
func (m *MySQL) cleanup() error {
	day := time.Now().Add(retention * -1).Format("01_02_2006")
	log.Info(fmt.Sprintf("cleaning up sip messages older than date [%s]", day))

	// drop the table
	_, err := m.db.Exec(fmt.Sprintf("DROP TABLE %s_%s", tablePrefix, day))
	if err != nil {
		return err
	}
	return nil
}
