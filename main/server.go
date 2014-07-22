package main

import (
	"hep"
	"net"
)

const hepId = 0x48455033

func UDPServer(ip string, port int) {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(ip),
	}
	conn, err := net.ListenUDP("udp", &addr)
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	for {
		packet := make([]byte, 2048)
		hepMsg := &hep.HepMsg{}

		_, _, err := conn.ReadFromUDP(packet)

		if err != nil {
			continue
		}

		hepMsg.Parse(packet)

		// Do something with the parsed message
	}
}
