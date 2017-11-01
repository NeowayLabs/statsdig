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

// Sampler abstraction, makes it easy to have multiple implementations
// of a sampler, which can be useful to testing
type Sampler interface {

	// Count sends a increment to a count metric with given name
	Count(name string, tags ...Tag) error

	// Gauge sets the gauge with the given name to the given value
	Gauge(name string, value int, tags ...Tag) error

	// Time sets the time in milliseconds with the given name to the given value
	Time(name string, value int, tags ...Tag) error
}

// UDPSampler is a sampler that sends metrics through UDP
type UDPSampler struct {
	conn net.PacketConn
	addr *net.UDPAddr
}

func (s *UDPSampler) write(data []byte) (int, error) {
	return s.conn.WriteTo(data, s.addr)
}

// NewSysdigSampler creates a sampler suited to work
// with the sysdig cloud client, sending metrics to localhost
// at the default statsd port.
func NewSysdigSampler() (*UDPSampler, error) {
	return NewSampler("127.0.0.1:8125")
}

// NewSampler creates a sampler suited to work
// with any statsd server listening add the given addr,
// where addr must be serializeted just as the addr provided
// to Go's net.ResolveUDPAddr function.
func NewSampler(addr string) (*UDPSampler, error) {
	udpAddr, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		return nil, fmt.Errorf("resolve udp address failed: %s", err)
	}
	conn, err := net.ListenPacket("udp4", ":0")
	if err != nil {
		return nil, fmt.Errorf("connection creation failed: %s", err)
	}
	return &UDPSampler{
		conn: conn,
		addr: udpAddr,
	}, nil
}

// Count sends a counter metric as specified here:
// https://github.com/b/statsd_spec#counters
func (sampler *UDPSampler) Count(name string, tags ...Tag) error {
	countType := "c"
	message := serialize(name, 1, countType, tags...)
	return sampler.send(message)
}

// Gauge sends a gauge metric as specified here:
// https://github.com/b/statsd_spec#gauges
func (sampler *UDPSampler) Gauge(name string, value int, tags ...Tag) error {
	countType := "g"
	message := serialize(name, value, countType, tags...)
	return sampler.send(message)
}

// Time sends a time metric as specified here:
// https://github.com/b/statsd_spec#timers
func (sampler *UDPSampler) Time(name string, value int, tags ...Tag) error {
	timeType := "ms"
	message := serialize(name, value, timeType, tags...)
	return sampler.send(message)
}
func (sampler *UDPSampler) send(message []byte) error {
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
		"%s%s:%d|%s",
		name,
		fulltags,
		value,
		metricType,
	))
}
