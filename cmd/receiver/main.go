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

type metric struct {
	addr net.Addr
	msg  string
}

func listener(port int, received chan<- metric) {
	conn, err := net.ListenPacket("udp", fmt.Sprintf(":%d", port))
	panicAtTheDisco(err)

	packet := make([]byte, 1024)
	for {
		_, addr, err := conn.ReadFrom(packet)
		panicAtTheDisco(err)
		received <- metric{addr: addr, msg: string(packet)}
	}
}

func main() {
	var port int
	flag.IntVar(&port, "port", 8125, "port to listen to")

	counter := 0
	received := make(chan metric)

	go listener(port, received)

	for {
		log.Printf("Listening for packages at: %d", port)
		m := <-received
		counter += 1
		log.Printf("Read: %s from: %s", m.msg, m.addr)
		log.Printf("Total received: %d", counter)
	}
}
