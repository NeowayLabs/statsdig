package statsdig

import (
	"fmt"
	"sync"
	"time"
)

// MemSampler is a sampler that gathers metrics in memory
// and allows the count of metrics to be queried.
// It's only usage is to test metric sampling on your application.
type MemSampler struct {
	sync.Mutex
	storage      map[string]int
	storageFloat map[string]float64
}

func NewMemSampler() *MemSampler {
	return &MemSampler{
		storage:      map[string]int{},
		storageFloat: map[string]float64{},
	}
}

func (s *MemSampler) Count(name string, tags ...Tag) error {
	serialized := serializeCount(name, tags)
	s.add(serialized)
	return nil
}

func (s *MemSampler) GetCount(name string, tags ...Tag) int {
	serialized := serializeCount(name, tags)
	return s.get(serialized)
}

func (s *MemSampler) Gauge(name string, value int, tags ...Tag) error {
	serialized := serializeGauge(name, value, tags)
	s.add(serialized)
	return nil
}

func (s *MemSampler) GetGauge(name string, value int, tags ...Tag) int {
	serialized := serializeGauge(name, value, tags)
	return s.get(serialized)
}

func (s *MemSampler) GaugeFloat(name string, value float64, tags ...Tag) error {
	serialized := serializeGaugeFloat(name, value, tags)
	s.add(serialized)
	return nil
}

func (s *MemSampler) GetGaugeFloat(name string, value float64, tags ...Tag) float64 {
	serialized := serializeGaugeFloat(name, value, tags)
	return s.getFloat(serialized)
}

func (s *MemSampler) Time(name string, value time.Duration, tags ...Tag) error {
	serialized := serializeTime(name, value, tags)
	s.add(serialized)
	return nil
}

func (s *MemSampler) GetTime(name string, value time.Duration, tags ...Tag) int {
	serialized := serializeTime(name, value, tags)
	return s.get(serialized)
}
func serializeCount(name string, tags []Tag) string {
	return fmt.Sprintf("count:%s:%v", name, tags)
}

func serializeGauge(name string, value int, tags []Tag) string {
	return fmt.Sprintf("gauge:%s:%d:%v", name, value, tags)
}

func serializeGaugeFloat(name string, value float64, tags []Tag) string {
	return fmt.Sprintf("gaugefloat:%s:%v:%v", name, value, tags)
}

func serializeTime(name string, value time.Duration, tags []Tag) string {
	return fmt.Sprintf("time:%s:%d:%v", name, value, tags)
}

func (s *MemSampler) add(serialized string) {
	s.Lock()
	s.storage[serialized] += 1
	s.storageFloat[serialized] += 1
	s.Unlock()
}

func (s *MemSampler) get(serialized string) int {
	s.Lock()
	defer s.Unlock()
	return s.storage[serialized]
}

func (s *MemSampler) getFloat(serialized string) float64 {
	s.Lock()
	defer s.Unlock()
	return s.storageFloat[serialized]
}
