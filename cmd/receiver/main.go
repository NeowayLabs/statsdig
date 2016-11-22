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

func listener(conn net.PacketConn, received chan<- metric) {

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
	received := make(chan metric, 10000)

	conn, err := net.ListenPacket("udp", fmt.Sprintf(":%d", port))
	panicAtTheDisco(err)

	listeners := 100

	for i := 0; i < listeners; i++ {
		go listener(conn, received)
	}

	for {
		log.Printf("Listening for packages at: %d", port)
		m := <-received
		counter += 1
		log.Printf("Read: %s from: %s", m.msg, m.addr)
		log.Printf("Total received: %d", counter)
	}
}
