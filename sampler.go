// statsdig is a minimalist statsd client focused
// on integrating with sysdig cloud services.
// Although it can work with plain StatsD just fine.
package statsdig

import (
	"fmt"

	"github.com/cactus/go-statsd-client/statsd"
)

// Sampler is how you will be able to sample metrics.
type Sampler struct {
	client statsd.Statter
}

// NewSysdigSampler creates a sampler suited to work
// with the sysdig cloud client, sending metrics to localhost
// at the default statsd port.
func NewSysdigSampler() (*Sampler, error) {
	return NewSampler("127.0.0.1:8125")
}

// NewSampler creates a sampler suited to work
// with any statsd server listening add the given addr,
// where addr must be formatted just as the addr provided
// to Go's net.ResolveUDPAddr function.
func NewSampler(addr string) (*Sampler, error) {
	client, err := statsd.NewClient(addr, "test-client")
	if err != nil {
		return nil, fmt.Errorf("resolve udp address failed: %s", err)
	}
	return &Sampler{
		client: client,
	}, nil
}

// Count sends a counter metric as specified here:
// https://github.com/b/statsd_spec#counters
func (sampler *Sampler) Count(name string) error {
	return sampler.client.Inc(name, 1, 1.0)
}
