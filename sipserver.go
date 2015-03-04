package flanders

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/weave-lab/flanders/db"
	"github.com/weave-lab/flanders/hep"
	"github.com/weave-lab/flanders/log"

	// Choose your db handler or import your own here
	// _ "lab.getweave.com/weave/flanders/db/influx"
	_ "github.com/weave-lab/flanders/db/mongo"
)

var h = ListenerHub{
	broadcast:   make(chan db.DbObject),
	register:    make(chan *Listener),
	unregister:  make(chan *Listener),
	connections: make(map[*Listener]bool),
}

func UDPServer(ip string, port int) {
	go h.run()
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(ip),
	}
	log.Info("Flanders server listening on " + ip + ":" + strconv.Itoa(port))
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Crit(err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	for {
		packet := make([]byte, 4096)

		length, _, err := conn.ReadFromUDP(packet)

		if err != nil {
			log.Err(err.Error())
			continue
		}

		packet = packet[:length]

		//log.Debug(string(packet))
		fmt.Printf("[% 3X]", packet)

		hepMsg, hepErr := hep.NewHepMsg(packet)

		if hepErr != nil {
			log.Err(hepErr.Error())
			continue
		}

		switch hepMsg.SipMsg.StartLine.Method {
		case "OPTIONS":
			continue
		case "SUBSCRIBE":
			continue
		case "NOTIFY":
			continue
			// case "REGISTER":
			// 	continue
		}

		switch hepMsg.SipMsg.Cseq.Method {
		case "OPTIONS":
			continue
		case "SUBSCRIBE":
			continue
		case "NOTIFY":
			continue
			// case "REGISTER":
			// 	continue
		}

		var datetime db.Time

		log.Debug(string(packet))
		if hepMsg.Timestamp != 0 {
			datetime = db.Time{time.Unix(int64(hepMsg.Timestamp), int64(hepMsg.TimestampMicro)*1000)}
		} else {
			datetime = db.Time{time.Now()}
		}

		dbObject := db.NewDbObject()
		dbObject.Datetime = datetime
		dbObject.MicroSeconds = datetime.Nanosecond() / 1000
		dbObject.Method = hepMsg.SipMsg.StartLine.Method + hepMsg.SipMsg.StartLine.Resp
		dbObject.ReplyReason = hepMsg.SipMsg.StartLine.RespText
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
		dbObject.UserAgent = hepMsg.SipMsg.UserAgent
		dbObject.Cseq = hepMsg.SipMsg.Cseq.Val
		for _, header := range hepMsg.SipMsg.Headers {
			if header.Header == "x-cid" {
				dbObject.CallIdAleg = header.Val
			}
		}

		dbObject.Msg = hepMsg.SipMsg.Msg

		h.broadcast <- *dbObject

		err = dbObject.Save()
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}

type Listener struct {
	Filter    string
	Broadcast chan db.DbObject
	quit      chan bool
}

type ListenerHub struct {
	// Registered connections.
	connections map[*Listener]bool

	// Inbound messages from the udp server
	broadcast chan db.DbObject

	// Register requests from the connections.
	register chan *Listener

	// Unregister requests from connections.
	unregister chan *Listener

	// Quit channel
	quit chan bool
}

func (h *ListenerHub) run() {
hubforloop:
	for {
		select {
		case l := <-h.register:
			h.connections[l] = true
		case l := <-h.unregister:
			if _, ok := h.connections[l]; ok {
				delete(h.connections, l)
			}
		case m := <-h.broadcast:
			for l := range h.connections {
				if strings.Contains(strings.ToLower(m.Msg), strings.ToLower(l.Filter)) {
					l.Broadcast <- m
				}
			}
		case <-h.quit:
			for _ = range h.connections {
				h.quit <- true
			}
			break hubforloop
		}
	}
}

func RegisterListener(filter string) Listener {
	newlistener := Listener{
		Filter:    filter,
		Broadcast: make(chan db.DbObject, 100),
	}

	h.register <- &newlistener
	return newlistener
}
