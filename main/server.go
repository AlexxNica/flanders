package main

import (
	"fmt"
	"net"
)

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
		fmt.Println("from", remote_addr, "got message:", string(buf), " Error: ", err, n)
	}
}
