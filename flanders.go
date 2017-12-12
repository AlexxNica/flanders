package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/weave-lab/flanders/api"
	"github.com/weave-lab/flanders/capture"
	"github.com/weave-lab/flanders/db"
	"github.com/weave-lab/flanders/log"

	// Import all DB handlers, use config to select which one
	// _ "lab.getweave.com/weave/flanders/db/influx"
	_ "github.com/weave-lab/flanders/db/mongo"
	_ "github.com/weave-lab/flanders/db/mysql"
)

func main() {
	webPort := flag.String("webport", "8000", "Web server port")
	sipPort := flag.String("sipport", "9060", "SIP server port")
	dbConnectString := flag.String("db", "localhost", "DB connect string")
	dbDriver := flag.String("driver", "mongo", "db driver (mysql|mongo|influx)")
	assetFolder := flag.String("assets", "public", "Static assets folder for GUI")

	loglevel := flag.String("loglevel", "warn", "Log level")
	flag.Parse()
	log.SetLogger(os.Stdout)
	log.SetLogLevel(*loglevel)
	log.Debug("debug enabled")

	webAddress := "0.0.0.0:" + *webPort
	sipAddress := "0.0.0.0:" + *sipPort

	err := db.Setup(*dbDriver, *dbConnectString)
	if err != nil {
		fmt.Printf("\nunable to connect to db: %s\nn", err)
		os.Exit(1)
	}

	err = capture.StartSIPCaptureServer(sipAddress)
	if err != nil {
		fmt.Printf("unable to start sip listener: %s", err)
		os.Exit(1)
	}

	// This blocks so do last
	err = api.StartWebServer(webAddress, *assetFolder)
	if err != nil {
		fmt.Printf("unable to start web server: %s", err)
		os.Exit(1)
	}
}
