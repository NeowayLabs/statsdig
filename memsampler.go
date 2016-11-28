package statsdig

import (
	"fmt"
	"sync"
)

// MemSampler is a sampler that gathers metrics in memory
// and allows the count of metrics to be queried.
// It's only usage is to test metric sampling on your application.
type MemSampler struct {
	sync.Mutex
	storage map[string]int
}

func NewMemSampler() *MemSampler {
	return &MemSampler{
		storage: map[string]int{},
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

func serializeCount(name string, tags []Tag) string {
	return fmt.Sprintf("count:%s:%v", name, tags)
}

func serializeGauge(name string, value int, tags []Tag) string {
	return fmt.Sprintf("gauge:%s:%d:%v", name, value, tags)
}

func (s *MemSampler) add(serialized string) {
	s.Lock()
	s.storage[serialized] += 1
	s.Unlock()
}

func (s *MemSampler) get(serialized string) int {
	s.Lock()
	defer s.Unlock()
	return s.storage[serialized]
}
