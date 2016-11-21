// sender generates lots of metrics so you can debug on your dash
package main

import (
	"flag"
	"github.com/NeowayLabs/statsdig"
	"log"
)

func panicAtTheDisco(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var count int
	var metric string

	flag.IntVar(&count, "count", 10000, "amount of counts to perform")
	flag.StringVar(&metric, "metric", "statsdig.test", "metric name to test")
	flag.Parse()

	log.Printf("Starting sampler metric[%s] count[%d]", metric, count)
	sampler, err := statsdig.NewSysdigSampler("127.0.0.1:8125")
	panicAtTheDisco(err)

	for i := 0; i < count; i++ {
		sampler.Count(metric)
	}

	log.Println("Done")
}
