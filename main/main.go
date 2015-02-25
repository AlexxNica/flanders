package main

import (
	"flag"
	"os"

	"github.com/weave-lab/flanders"
	"github.com/weave-lab/flanders/log"
)

func main() {
	flag.Parse()
	log.SetLogger(os.Stdout)
	log.SetLogLevel("debug")

	go flanders.UDPServer("0.0.0.0", 9060)
	flanders.WebServer("0.0.0.0", 8080)
	// quit := make(chan struct{})
	// <-quit
}
