// sender generates lots of metrics so you can debug on your dash
package main

import (
	"flag"
	"log"
	"time"

	"github.com/NeowayLabs/statsdig"
)

func panicAtTheDisco(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var count int
	var metric string
	delay := 5 * time.Millisecond

	flag.IntVar(&count, "count", 100000, "amount of counts to perform")
	flag.DurationVar(&delay, "delay", delay, "delay in ms before sending metrics")
	flag.StringVar(&metric, "metric", "statsdig.test", "metric name to test")
	flag.Parse()

	log.Printf("Starting sampler metric[%s] count[%d]", metric, count)
	sampler, err := statsdig.NewSysdigSampler()
	panicAtTheDisco(err)

	for i := 0; i < count; i++ {
		err := sampler.Count(metric)
		panicAtTheDisco(err)
		time.Sleep(delay)
	}

	log.Println("Done")
}
