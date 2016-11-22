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

	counter := 0
	conn, err := net.ListenPacket("udp", fmt.Sprintf(":%d", port))
	panicAtTheDisco(err)

	packet := make([]byte, 1024)
	for {
		_, _, err := conn.ReadFrom(packet)
		panicAtTheDisco(err)
		counter += 1
		log.Println(string(packet), counter)
	}
}
