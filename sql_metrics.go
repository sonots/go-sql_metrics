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
var proxyRegistry = make(map[string](*DB))

//Wrap  instrument template
func WrapDB(name string, db *sql.DB) *DB {
	proxy := newDB(name, db)
	proxyRegistry[name] = proxy
	return proxy
}

//Print  print the metrics in each second
func Print(duration int) {
	timeDuration := time.Duration(duration)
	go func() {
		time.Sleep(timeDuration * time.Second)
		for {
			startTime := time.Now()
			for _, proxy := range proxyRegistry {
				proxy.printMetrics(duration)
			}
			elapsedTime := time.Now().Sub(startTime)
			time.Sleep(timeDuration*time.Second - elapsedTime)
		}
	}()
}
