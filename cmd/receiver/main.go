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

	counter := 0
	packet := make([]byte, 65536)

	for {
		log.Printf("Listening for packages at: %d", port)
		_, _, err := conn.ReadFrom(packet)
		panicAtTheDisco(err)
		counter += 1
		log.Printf("Read: %d from: %s", n, addr)
		log.Printf("Metric:'%s'", string(packet))
		log.Printf("Total received: %d", counter)
	}
}
