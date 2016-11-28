// statsdig is a minimalist statsd client focused
// on integrating with sysdig cloud services.
// Although it can work with plain StatsD just fine.
package statsdig

import (
	"fmt"
	"net"
	"strings"
)

// Tag represents a Sysdig StatsD tag, which is a extension
// to add more dimensions to your metric, like prometheus labels.
// More: https://support.sysdigcloud.com/hc/en-us/articles/204376099-Metrics-integrations-StatsD
type Tag struct {
	Name  string
	Value string
}

// Sampler abstraction, makes it easy to test metric generation
type Sampler interface {

	// Send a 1 increment to a count metric with given name
	Count(name string, tags ...Tag) error
}

// sampler is how you will be able to sample metrics.
type statsd struct {
	conn net.PacketConn
	addr *net.UDPAddr
}

func (s *statsd) write(data []byte) (int, error) {
	return s.conn.WriteTo(data, s.addr)
}

// NewSysdigSampler creates a sampler suited to work
// with the sysdig cloud client, sending metrics to localhost
// at the default statsd port.
func NewSysdigSampler() (*statsd, error) {
	return NewSampler("127.0.0.1:8125")
}

// NewSampler creates a sampler suited to work
// with any statsd server listening add the given addr,
// where addr must be serializeted just as the addr provided
// to Go's net.ResolveUDPAddr function.
func NewSampler(addr string) (*statsd, error) {
	udpAddr, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		return nil, fmt.Errorf("resolve udp address failed: %s", err)
	}
	conn, err := net.ListenPacket("udp4", ":0")
	if err != nil {
		return nil, fmt.Errorf("connection creation failed: %s", err)
	}
	return &statsd{
		conn: conn,
		addr: udpAddr,
	}, nil
}

// Count sends a counter metric as specified here:
// https://github.com/b/statsd_spec#counters
func (sampler *statsd) Count(name string, tags ...Tag) error {
	countType := "c"
	message := serialize(name, 1, countType, tags...)
	n, err := sampler.write(message)
	if err != nil {
		return err
	}
	if n != len(message) {
		return fmt.Errorf(
			"expected to send %d but sent %d",
			len(message),
			n,
		)
	}
	return nil
}

func serialize(
	name string,
	value int,
	metricType string,
	tags ...Tag,
) []byte {
	var strtags []string
	for _, tag := range tags {
		strtags = append(strtags, tag.Name+"="+tag.Value)
	}
	fulltags := ""
	if len(strtags) > 0 {
		fulltags += "#" + strings.Join(strtags, ",")
	}
	return []byte(fmt.Sprintf(
		"%s%s:1|c",
		name,
		fulltags,
	))
}
