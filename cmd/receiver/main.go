package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

func panicAtTheDisco(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var port int
	flag.IntVar(&port, "port", 8125, "port to listen to")

	conn, err := net.ListenPacket("udp", fmt.Sprintf(":%d", port))
	panicAtTheDisco(err)

	for {
		packet := make([]byte, 65536)
		log.Printf("Listening for packages at: %d", port)
		n, addr, err := conn.ReadFrom(packet)
		panicAtTheDisco(err)
		log.Printf("Read: %d from: %s", n, addr)
		log.Printf("Packet: %s", string(packet))
	}
}
