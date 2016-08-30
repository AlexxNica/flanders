package flanders

import (
	"fmt"

	"github.com/weave-lab/flanders/db"

	// Choose your db handler or import your own here
	// _ "lab.getweave.com/weave/flanders/db/influx"
	_ "github.com/weave-lab/flanders/db/mongo"
)

func Run(sipAddress string, webAddress string, dbAddress string) error {

	err := db.Setup("mongo", dbAddress)
	if err != nil {
		return fmt.Errorf("unable to connect to db: %s", err)
	}

	go h.run()

	err = StartSIPServer(sipAddress)
	if err != nil {
		return fmt.Errorf("unable to start sip listener: %s", err)
	}

	err = StartWebServer(webAddress)
	if err != nil {
		return fmt.Errorf("unable to start web server: %s", err)
	}

	select {}
}
