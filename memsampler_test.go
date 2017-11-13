package statsdig_test

import (
	"sync"
	"testing"
	"time"

	"github.com/NeowayLabs/statsdig"
)

func TestMemSampler(t *testing.T) {
	count := 100
	var wg sync.WaitGroup
	s := statsdig.NewMemSampler()
	tags := []statsdig.Tag{
		statsdig.Tag{
			Name:  "mem",
			Value: "1",
		},
	}

	wg.Add(count)
	expectedTime := time.Duration(1000 * time.Millisecond)
	for i := 0; i < count; i++ {
		go func() {
			s.Count("count")
			s.Count("count", tags...)
			s.Gauge("gauge", 666)
			s.Gauge("gauge", 777, tags...)
			s.Time("time", expectedTime)
			s.Time("time", expectedTime, tags...)
			wg.Done()
		}()
	}
	wg.Wait()

	checkMetric := func(name string, counter func() int) {
		if counter() != count {
			t.Fatalf(
				"%s: expected %d but got %d",
				name,
				counter(),
				count,
			)
		}
	}

	checkMetric("count", func() int {
		return s.GetCount("count")
	})

	checkMetric("countWithTags", func() int {
		return s.GetCount("count", tags...)
	})

	checkMetric("gauge", func() int {
		return s.GetGauge("gauge", 666)
	})

	checkMetric("gaugeWithTags", func() int {
		return s.GetGauge("gauge", 777, tags...)
	})

	checkMetric("time", func() int {
		return s.GetTime("time", expectedTime)
	})

	checkMetric("timeWithTags", func() int {
		return s.GetTime("time", expectedTime, tags...)
	})
}
