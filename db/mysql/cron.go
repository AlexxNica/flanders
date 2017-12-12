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

//remove old records daily at 7:00 a.m. UTC (1:00 a.m. MST)
func (m *MySQL) startCron() {
	log.Info(fmt.Sprintf("setup cron to delete messages older than %d days", *cleanOldExpired))
	retention = time.Hour * 24 * time.Duration(*cleanOldExpired)
	c := cron.New()
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
	since := time.Now().Add(retention * -1).Format(time.RFC3339)
	log.Info(fmt.Sprintf("cleaning up sip messages older than date [%s]", since))
	_, err := m.db.Exec("DELETE FROM messages WHERE date < ?", since)
	if err != nil {
		return err
	}
	return nil
}
