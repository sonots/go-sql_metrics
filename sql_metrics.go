package sql_metrics

import (
	"database/sql"
	"time"
)

// print infomation on each request
var Verbose = false

// Set Enable = false if you want to turn off the instrumentation
var Enable = true

// a set of proxies
var proxyRegistry = make(map[*sql.DB](*DB))

// a set of metrics
var metricsRegistry = make(map[string](*Metrics))

//WrapDB
func WrapDB(name string, db *sql.DB) *DB {
	metrics := metricsRegistry[name]
	if metrics == nil {
		metrics = newMetrics(name)
		metricsRegistry[name] = metrics
	}
	proxy := proxyRegistry[db]
	if proxy == nil {
		proxy = newDB(db, metrics)
		proxyRegistry[db] = proxy
	}
	return proxy
}

//Flush  print the metrics
func Flush() {
	for _, metrics := range metricsRegistry {
		metrics.printMetrics(-1)
	}
}

//Print  print the metrics in each second
func Print(duration int) {
	timeDuration := time.Duration(duration)
	go func() {
		time.Sleep(timeDuration * time.Second)
		for {
			startTime := time.Now()
			for _, metrics := range metricsRegistry {
				metrics.printMetrics(duration)
			}
			elapsedTime := time.Now().Sub(startTime)
			time.Sleep(timeDuration*time.Second - elapsedTime)
		}
	}()
}
