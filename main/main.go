package main

import (
	"flag"

	"github.com/weave-lab/flanders"
	"github.com/weave-lab/flanders/log"
)

func main() {
	//webport := flag.Int("webport", 80, "Web server port")
	sipport := flag.Int("sipport", 9060, "SIP server port")
	loglevel := flag.String("loglevel", "warn", "Log level")
	flag.Parse()
	//log.SetLogger(os.Stdout)
	log.SetLogLevel(*loglevel)

	flanders.UDPServer("0.0.0.0", *sipport)
	//flanders.WebServer("0.0.0.0", *webport)
	// quit := make(chan struct{})
	// <-quit
}
