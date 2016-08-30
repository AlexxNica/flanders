package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/weave-lab/flanders"
	"github.com/weave-lab/flanders/log"
)

func main() {
	webPort := flag.String("webport", "8000", "Web server port")
	sipPort := flag.String("sipport", "9060", "SIP server port")
	dbAddress := flag.String("db", "localhost", "DB address")

	loglevel := flag.String("loglevel", "warn", "Log level")
	flag.Parse()
	//log.SetLogger(os.Stdout)
	log.SetLogLevel(*loglevel)

	webAddress := "0.0.0.0:" + *webPort
	sipAddress := "0.0.0.0:" + *sipPort

	err := flanders.Run(sipAddress, webAddress, *dbAddress)
	if err != nil {
		fmt.Println(err)
		log.Crit(err.Error())
		os.Exit(1)
	}
}
