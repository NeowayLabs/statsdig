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
	result string,
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

	for i := 0; i < count; i++ {
		msg := listener.Get(i, timeout)
		if msg != result {
			t.Fatalf("Expected %q but got %q", result, msg)
		}
	}
	msg := listener.Get(count, 100*time.Millisecond)
	if msg != "" {
		t.Fatalf("Received unexpected msg: %s", msg)
	}
}

type countcase struct {
	Name   string
	Metric string
	Tags   []statsdig.Tag
	Result string
}

type gaugecase struct {
	Name   string
	Metric string
	Tags   []statsdig.Tag
	Result string
	Gauge  int
}

type gaugefloatcase struct {
	Name       string
	Metric     string
	Tags       []statsdig.Tag
	Result     string
	GaugeFloat float64
}

func ExampleUDPSampler_Count() {
	// Creating a Sysdig specific sampler
	sampler, err := statsdig.NewSysdigSampler()
	if err != nil {
		return
	}
	sampler.Count("metric.name", statsdig.Tag{
		Name:  "tagname",
		Value: "tagvalue",
	})
}

func TestCount(t *testing.T) {

	cases := []countcase{
		countcase{
			Name:   "testCount",
			Metric: "TestCount",
			Result: "TestCount:1|c",
		},
		countcase{
			Name:   "testCountWithTag",
			Metric: "TestCountTag",
			Tags: []statsdig.Tag{
				statsdig.Tag{
					Name:  "tag",
					Value: "hi",
				},
			},
			Result: "TestCountTag#tag=hi:1|c",
		},
		countcase{
			Name:   "testCountWithTags",
			Metric: "TestCountTags",
			Tags: []statsdig.Tag{
				statsdig.Tag{
					Name:  "tag",
					Value: "hi",
				},
				statsdig.Tag{
					Name:  "tag2",
					Value: "1",
				},
			},
			Result: "TestCountTags#tag=hi,tag2=1:1|c",
		},
	}

	for _, tcase := range cases {
		t.Run(tcase.Name, func(t *testing.T) {
			testMetric(
				t,
				func(t *testing.T, sampler statsdig.Sampler) error {
					return sampler.Count(tcase.Metric, tcase.Tags...)
				},
				tcase.Result,
			)
		})
	}
}

func TestGauge(t *testing.T) {

	cases := []gaugecase{
		gaugecase{
			Name:   "testGauge",
			Metric: "TestGauge",
			Gauge:  500,
			Result: "TestGauge:500|g",
		},
		gaugecase{
			Name:   "testGaugeWithTag",
			Metric: "TestGaugeTag",
			Gauge:  666,
			Tags: []statsdig.Tag{
				statsdig.Tag{
					Name:  "tag",
					Value: "gauging",
				},
			},
			Result: "TestGaugeTag#tag=gauging:666|g",
		},
		gaugecase{
			Name:   "testGaugeWithTags",
			Metric: "TestGaugeTags",
			Gauge:  10,
			Tags: []statsdig.Tag{
				statsdig.Tag{
					Name:  "tag",
					Value: "hi",
				},
				statsdig.Tag{
					Name:  "tag2",
					Value: "1",
				},
			},
			Result: "TestGaugeTags#tag=hi,tag2=1:10|g",
		},
	}

	for _, tcase := range cases {
		t.Run(tcase.Name, func(t *testing.T) {
			testMetric(
				t,
				func(t *testing.T, sampler statsdig.Sampler) error {
					return sampler.Gauge(
						tcase.Metric,
						tcase.Gauge,
						tcase.Tags...,
					)
				},
				tcase.Result,
			)
		})
	}
}

func TestGaugeFloat(t *testing.T) {

	cases := []gaugefloatcase{
		gaugefloatcase{
			Name:       "testGaugeFloat",
			Metric:     "testGaugeFloat",
			GaugeFloat: 500.012,
			Result:     "testGaugeFloat:500.012|gf",
		},
		gaugefloatcase{
			Name:       "testGaugeFloatWithTag",
			Metric:     "TestGaugeFloatTag",
			GaugeFloat: 666.12,
			Tags: []statsdig.Tag{
				statsdig.Tag{
					Name:  "tag",
					Value: "gaugingfloat",
				},
			},
			Result: "TestGaugeFloatTag#tag=gaugingfloat:666.12|gf",
		},
		gaugefloatcase{
			Name:       "testGaugeFloatWithTags",
			Metric:     "TestGaugeFloatTags",
			GaugeFloat: 10.155,
			Tags: []statsdig.Tag{
				statsdig.Tag{
					Name:  "tag",
					Value: "hi",
				},
				statsdig.Tag{
					Name:  "tag2",
					Value: "1",
				},
			},
			Result: "TestGaugeFloatTags#tag=hi,tag2=1:10.155|gf",
		},
	}

	for _, tcase := range cases {
		t.Run(tcase.Name, func(t *testing.T) {
			testMetric(
				t,
				func(t *testing.T, sampler statsdig.Sampler) error {
					return sampler.GaugeFloat(
						tcase.Metric,
						tcase.GaugeFloat,
						tcase.Tags...,
					)
				},
				tcase.Result,
			)
		})
	}
}

func TestTime(t *testing.T) {

	type testcase struct {
		Name   string
		Metric string
		Tags   []statsdig.Tag
		Result string
		Value  time.Duration
	}

	cases := []testcase{
		testcase{
			Name:   "testTime",
			Metric: "testTime",
			Value:  time.Duration(1 * time.Millisecond),
			Result: "testTime:1|ms",
		},
		testcase{
			Name:   "testTimeWithTag",
			Metric: "testTimeTag",
			Value:  time.Duration(1 * time.Minute),
			Tags: []statsdig.Tag{
				statsdig.Tag{
					Name:  "tag",
					Value: "hi",
				},
			},
			Result: "testTimeTag#tag=hi:60000|ms",
		},
		testcase{
			Name:   "testTimeWithTags",
			Metric: "testTimeTags",
			Value:  time.Duration(1000000),
			Tags: []statsdig.Tag{
				statsdig.Tag{
					Name:  "tag",
					Value: "hi",
				},
				statsdig.Tag{
					Name:  "tag2",
					Value: "1",
				},
			},
			Result: "testTimeTags#tag=hi,tag2=1:1|ms",
		},
	}

	for _, tcase := range cases {
		t.Run(tcase.Name, func(t *testing.T) {
			testMetric(
				t,
				func(t *testing.T, sampler statsdig.Sampler) error {
					return sampler.Time(tcase.Metric, tcase.Value, tcase.Tags...)
				},
				tcase.Result,
			)
		})
	}
}
