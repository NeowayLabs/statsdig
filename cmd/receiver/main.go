package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"sync/atomic"
)

func panicAtTheDisco(err error) {
	if err != nil {
		panic(err)
	}
}

var counter uint64

func listener(conn net.PacketConn) {

	packet := make([]byte, 1024)
	for {
		_, _, err := conn.ReadFrom(packet)
		panicAtTheDisco(err)

		atomic.AddUint64(&counter, 1)

	}
}

func main() {
	var port int
	flag.IntVar(&port, "port", 8125, "port to listen to")

	conn, err := net.ListenPacket("udp", fmt.Sprintf(":%d", port))
	panicAtTheDisco(err)

	go listener(conn)
	go listener(conn)
	go listener(conn)
	go listener(conn)
	go listener(conn)

	log.Printf("Listening for packages at: %d", port)

	for {
		time.Sleep(3 * time.Second)

		log.Printf("Total received: %d", atomic.LoadUint64(&counter))
	}
}
