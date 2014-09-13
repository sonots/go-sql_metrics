package sql_metrics

import (
	"fmt"
	metrics "github.com/sonots/go-metrics" // max,mean,min,stddev,percentile
	"time"
)

type Metrics struct {
	name   string
	timers map[string]metrics.Timer
}

func newMetrics(name string) *Metrics {
	return &Metrics{
		name:   name,
		timers: map[string]metrics.Timer{},
	}
}

//print the elapsed time on each request if Verbose flag is true
func (proxy *Metrics) printVerbose(elapsedTime time.Duration, query string) {
	fmt.Printf("time:%v\tdb:%s\tquery:%s\telapsed:%f\n",
		time.Now(),
		proxy.name,
		query,
		elapsedTime.Seconds(),
	)
}

func (proxy *Metrics) printMetrics(duration int) {
	for query, timer := range proxy.timers {
		count := timer.Count()
		if count > 0 {
			fmt.Printf(
				"time:%v\tdb:%s\tquery:%s\tcount:%d\tmax:%f\tmean:%f\tmin:%f\tpercentile95:%f\tsum:%f\tduration:%d\n",
				time.Now(),
				proxy.name,
				query,
				timer.Count(),
				float64(timer.Max())/float64(time.Second),
				timer.Mean()/float64(time.Second),
				float64(timer.Min())/float64(time.Second),
				timer.Percentile(0.95)/float64(time.Second),
				float64(timer.Sum())/float64(time.Second),
				duration,
			)
			proxy.timers[query] = metrics.NewTimer()
		}
	}
}

//measure the time
func (proxy *Metrics) measure(startTime time.Time, query string) {
	elapsedTime := time.Now().Sub(startTime)
	if proxy.timers[query] == nil {
		proxy.timers[query] = metrics.NewTimer()
	}
	proxy.timers[query].Update(elapsedTime)
	if Enable && Verbose {
		proxy.printVerbose(elapsedTime, query)
	}
}
