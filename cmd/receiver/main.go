package main

import (
	"flag"
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

func panicAtTheDisco(err error) {
	if err != nil {
		panic(err)
	}
}

func listener(conn net.PacketConn, counter *int64) {

	packet := make([]byte, 1024)
	for {
		_, _, err := conn.ReadFrom(packet)
		panicAtTheDisco(err)
		atomic.AddInt64(counter, 1)
	}
}

func main() {
	var port int
	flag.IntVar(&port, "port", 8125, "port to listen to")

	var counter int64

	conn, err := net.ListenPacket("udp", fmt.Sprintf(":%d", port))
	panicAtTheDisco(err)

	listeners := 10

	for i := 0; i < listeners; i++ {
		go listener(conn, &counter)
	}

	for {
		time.Sleep(1 * time.Second)
		fmt.Println(counter)
	}
}
