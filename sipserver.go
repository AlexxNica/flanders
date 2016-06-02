package flanders

import (
	"fmt"
	"net"
	"os"
	"strings"

	"lab.getweave.com/weave/flanders/db"
	"lab.getweave.com/weave/flanders/log"
)

var h = ListenerHub{
	broadcast:   make(chan db.DbObject),
	register:    make(chan *Listener),
	unregister:  make(chan *Listener),
	connections: make(map[*Listener]bool),
}

func StartSIPServer(address string) error {

	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}

	log.Info("Flanders server listening on " + addr.String())
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Crit(err.Error())
		os.Exit(1)
	}

	go func() {
		for {
			packet := make([]byte, 4096)

			length, _, err := conn.ReadFromUDP(packet)
			if err != nil {
				log.Err(err.Error())
				continue
			}

			packet = packet[:length]

			err = processPacket(packet)
			if err != nil {
				log.Err(fmt.Sprintf("Unable to process packet: %s", err))
				continue
			}
		}

		conn.Close()

	}()

	return nil
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
