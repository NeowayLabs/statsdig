package statsdig_test

import (
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/NeowayLabs/statsdig"
)

type msgListener struct {
	sync.Mutex
	conn     net.PacketConn
	msgs     []string
	closed   chan struct{}
	isclosed bool
}

var port int = 8124

func getport() int {
	port += 1
	return port
}

func abortOnErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func (l *msgListener) Listen(t *testing.T, port int) {
	if l.conn != nil {
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
					n, _, err := l.conn.ReadFrom(packet)
					if err != nil {
						// May happen when we call close
						return
					}
					if n == 0 {
						continue
					}
					l.addMsg(string(packet[:n]))
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

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		l.Lock()
		ret := ""
		if i < len(l.msgs) {
			ret = l.msgs[i]
		}
		l.Unlock()
		if ret != "" {
			return ret
		}
		time.Sleep(10 * time.Millisecond)
	}
	return ""
}

func getlocalhost(port int) string {
	return fmt.Sprintf("127.0.0.1:%d", port)
}

func newListener() *msgListener {
	return &msgListener{
		isclosed: true,
	}
}

type samplerFunc func(*testing.T, statsdig.Sampler) error
type getExpectedMetricFunc func() string

func testMetric(
	t *testing.T,
	sample samplerFunc,
	getExpectedMetric getExpectedMetricFunc,
) {
	listener := newListener()
	defer listener.Close(t)

	port := getport()
	listener.Listen(t, port)
	sampler, err := statsdig.NewSampler(getlocalhost(port))
	abortOnErr(t, err)

	count := 10

	for i := 0; i < count; i++ {
		err := sample(t, sampler)
		abortOnErr(t, err)
	}

	timeout := 1 * time.Second
	expectedMetric := getExpectedMetric()

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

type testcase struct {
	name   string
	metric string
	tags   []statsdig.Tag
	result string
}

func TestCount(t *testing.T) {

	cases := []testcase{
		testcase{
			name:   "testCount",
			metric: "TestCount",
			result: "TestCount:1|c",
		},
		testcase{
			name:   "testCountWithTag",
			metric: "TestCountTag",
			tags: []statsdig.Tag{
				statsdig.Tag{
					Name:  "tag",
					Value: "hi",
				},
			},
			result: "TestCountTag#tag=hi:1|c",
		},
		testcase{
			name:   "testCountWithTags",
			metric: "TestCountTags",
			tags: []statsdig.Tag{
				statsdig.Tag{
					Name:  "tag",
					Value: "hi",
				},
				statsdig.Tag{
					Name:  "tag2",
					Value: "1",
				},
			},
			result: "TestCountTags#tag=hi,tag2=1:1|c",
		},
	}

	for _, tcase := range cases {
		t.Run(tcase.name, func(t *testing.T) {
			testMetric(
				t,
				func(t *testing.T, sampler statsdig.Sampler) error {
					return sampler.Count(tcase.metric, tcase.tags...)
				},
				func() string {
					return tcase.result
				},
			)
		})
	}
}
