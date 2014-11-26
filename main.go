package main

import (
	"encoding/json"
	"fmt"
	"github.com/spacemonkeygo/spacelog"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"lab.getweave.com/weave/flanders/db"
	_ "lab.getweave.com/weave/flanders/db/influx"
	"lab.getweave.com/weave/flanders/hep"
	_ "log"
	"net"
	"net/http"
)

func main() {
	logger := spacelog.GetLogger()
	logger.Debug("Testing logger")
	go UDPServer("0.0.0.0", 9060)
	WebServer("0.0.0.0", 8080)
	// quit := make(chan struct{})
	// <-quit
}

var test int

func WebServer(ip string, port int) {

	goji.Get("/", func(c web.C, w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	})

	goji.Get("/search", func(c web.C, w http.ResponseWriter, r *http.Request) {
		filter := db.Filter{}
		options := &db.Options{}

		r.ParseForm()
		order := r.Form["orderby"]
		options.Sort = order

		var results []db.DbObject

		db.Db.Find(&filter, options, &results)
		jsonResults, err := json.Marshal(results)
		if err != nil {
			fmt.Fprint(w, err)
			return
		}

		fmt.Fprintf(w, "%s", string(jsonResults))
	})

	goji.Get("/call/:id", func(c web.C, w http.ResponseWriter, r *http.Request) {
		callId := c.URLParams["id"]

	})

	goji.Serve()
}

func UDPServer(ip string, port int) {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(ip),
	}
	fmt.Println("Flanders server listening on ", ip+":", port)
	conn, err := net.ListenUDP("udp", &addr)
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	for {
		packet := make([]byte, 4096)

		length, _, err := conn.ReadFromUDP(packet)

		packet = packet[:length]
		// hepString := string(packet[:length])

		// fmt.Printf("\nPacket: %X\n", truncatedPacket)
		// fmt.Printf("\nPacket: %s\n", hepString)

		if err != nil {
			fmt.Println(err)
			continue
		}

		hepMsg, hepErr := hep.NewHepMsg(packet)

		if hepErr != nil {
			fmt.Println("ERROR PARSING HEP MESSAGE.................")
			fmt.Println(hepErr)
			continue
		}
		//fmt.Printf("%#v\n", hepMsg)

		// Store HEP message in database
		dbObject := db.NewDbObject()
		dbObject.Method = hepMsg.SipMsg.StartLine.Method + hepMsg.SipMsg.StartLine.RespText
		dbObject.SourceIp = hepMsg.Ip4SourceAddress
		dbObject.SourcePort = hepMsg.SourcePort
		dbObject.DestinationIp = hepMsg.Ip4DestinationAddress
		dbObject.DestinationPort = hepMsg.DestinationPort
		dbObject.CallId = hepMsg.SipMsg.CallId
		dbObject.FromUser = hepMsg.SipMsg.From.URI.User
		dbObject.FromDomain = hepMsg.SipMsg.From.URI.Host
		dbObject.FromTag = hepMsg.SipMsg.From.Tag
		dbObject.ToUser = hepMsg.SipMsg.To.URI.User
		dbObject.ToDomain = hepMsg.SipMsg.To.URI.Host
		dbObject.ToTag = hepMsg.SipMsg.To.Tag

		// dbObject.ContactUser = hepMsg.SipMsg.Contact.URI.User
		// dbOjbect.ContactIp =
		// dbOjbect.ContactPort =

		dbObject.Msg = hepMsg.SipMsg.Msg

		fmt.Printf("\n\nDbObject-----------\n%+v\n", dbObject)

		err = dbObject.Save()
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}
