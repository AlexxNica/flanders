package main

import (
	"fmt"
	"lab.getweave.com/weave/flanders/sip"
	"net"
)

func main() {
	UDPServer("127.0.0.1", 9060)
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

		_, _, err := conn.ReadFromUDP(packet)

		fmt.Printf("Packet: %X\n", packet)

		if err != nil {
			fmt.Println(err)
			continue
		}

		sipMsg, sipErr := sip.NewSipMsg(packet)

		if sipErr != nil {
			fmt.Println(sipErr)
			continue
		}
		fmt.Printf("%+v\n", sipMsg)

		// Do something with the parsed message
	}
}
