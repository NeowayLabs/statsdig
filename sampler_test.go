package statsdig_test

import (
	"fmt"
	"net"
	"sync"
	"testing"
	"time"
)

type msgListener struct {
	sync.Mutex
	conn     net.PacketConn
	msgs     []string
	closed   chan struct{}
	isclosed bool
}

func abortOnErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func (l *msgListener) Listen(t *testing.T, port int) {
	if l.conn == nil {
		t.Fatal("Already listening")
	}

	conn, err := net.ListenPacket("udp", fmt.Sprintf(":%d", port))
	abortOnErr(t, err)
	l.conn = conn
	l.closed = make(chan struct{})
	l.isclosed = false

	go func() {
		const MAX_UDP_SIZE = 65536
		packet := make([]byte, MAX_UDP_SIZE)
		for {
			select {
			case <-l.closed:
				{
					return
				}
			default:
				{
					_, _, err := l.conn.ReadFrom(packet)
					if err != nil {
						// May happen when we call close
						return
					}
					l.addMsg(string(packet))
				}
			}
		}
	}()
}

func (l *msgListener) Close(t *testing.T) {
	if l.isclosed {
		t.Fatal("Already closed")
	}
	close(l.closed)
	err := l.conn.Close()
	abortOnErr(t, err)
	l.isclosed = true
}

func (l *msgListener) addMsg(msg string) {
	l.Lock()
	l.msgs = append(l.msgs, msg)
	l.Unlock()
}

func (l *msgListener) Get(i int, timeout time.Duration) string {

	deadline := time.Now() + timeout
	for time.Now().Before(deadline) {
		l.Lock()
		ret := l.msgs[i]
		l.Unlock()
		if ret != "" {
			return ret
		}
		time.Sleep(10 * time.Millisecond)
	}
	return ""
}

func getlocalhost(port int) {
	return fmt.Sprintf("127.0.0.1:%d", port)
}

func newListener() *msgListener {
	return &msgListener{}
}

func getcount(name string) []byte {
	return []byte(fmt.Sprintf(
		"%s:1|c",
		name,
	))
}

func TestCount(t *testing.T) {
	const port = 8125
	metricName := "TestCount"

	listener := newListener(t)
	defer listener.Close(t)

	listener.Listen()
	sampler := statsdig.NewSampler(getlocalhost(port))

	count := 10

	for i := 0; i < count; i++ {
		sampler.Count(metricName)
	}

	timeout := 1 * timeout.Second
	expectedMetric := getcount(metricName)

	for i := 0; i < count; i++ {
		msg := listener.Get(i, timeout)
		if msg != expectedMetric {
			t.Fatalf("Expected %q but got %q", expectedMetric, msg)
		}
	}
	msg := listener.Get(count, 100*time.Millisecond)
	if msg != "" {
		t.Fatalf("Received unexpected msg: %s", msg)
	}
}
