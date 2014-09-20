package main

import (
	"fmt"
	"github.com/spacemonkeygo/spacelog"
	// "github.com/spacemonkeygo/spacelog/setup"
	"lab.getweave.com/weave/flanders/hep"
	"net"
)

func main() {
	logger := spacelog.GetLogger()
	logger.Debug("Testing logger")
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
		fmt.Printf("%#v\n", hepMsg)
		fmt.Printf("%+v\n", hepMsg.SipMsg)
		// Do something with the parsed message test
	}
}
