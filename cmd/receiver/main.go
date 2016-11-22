package main

import (
	"flag"
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

var counter uint64

func panicAtTheDisco(err error) {
	if err != nil {
		panic(err)
	}
}

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

	listeners := 5

	for i := 0; i < listeners; i++ {
		go listener(conn)
	}

	for {
		time.Sleep(1 * time.Second)
		fmt.Println(atomic.LoadUint64(&counter))
	}
}
