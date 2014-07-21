package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

const hepId = 0x48455033

func UDPServer() {
	addr := net.UDPAddr{
		Port: 9060,
		IP:   net.ParseIP("127.0.0.1"),
	}
	conn, err := net.ListenUDP("udp", &addr)
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	for {
		buf := make([]byte, 512)

		n, remote_addr, err := conn.ReadFromUDP(buf)
		if binary.BigEndian.Uint32(buf[:3]) != hepId {
			continue
		}

		length := binary.BigEndian.Uint32(buf[4:5])

		fmt.Println("from", remote_addr, "got message:", string(buf), " Error: ", err, n)
	}
}
