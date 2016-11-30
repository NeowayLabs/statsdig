package statsdig_test

import (
	"sync"
	"testing"

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
	for i := 0; i < count; i++ {
		go func() {
			s.Count("count")
			s.Count("count", tags...)
			s.Gauge("gauge", 666)
			s.Gauge("gauge", 777, tags...)
			wg.Done()
		}()
	}
	wg.Wait()

	checkCount := func(name string, counter func() int) {
		if counter() != count {
			t.Fatalf(
				"%s: expected %d but got %d",
				name,
				counter(),
				count,
			)
		}
	}

	checkCount("count", func() int {
		return s.GetCount("count")
	})

	checkCount("countWithTags", func() int {
		return s.GetCount("count", tags...)
	})

	checkCount("gauge", func() int {
		return s.GetGauge("gauge", 666)
	})

	checkCount("gaugeWithTags", func() int {
		return s.GetGauge("gauge", 777, tags...)
	})
}
