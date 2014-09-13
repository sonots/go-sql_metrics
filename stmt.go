package sql_metrics

import (
	"database/sql"
	"time"
)

type Stmt struct {
	proxyDB  *DB
	query    string
	Original *sql.Stmt
}

func newStmt(proxy *DB, stmt *sql.Stmt, query string) *Stmt {
	return &Stmt{
		proxyDB:  proxy,
		Original: stmt,
		query:    query,
	}
}

func (proxy *Stmt) measure(startTime time.Time) {
	proxy.proxyDB.measure(startTime, proxy.query)
}

// instrument Exec
func (proxy *Stmt) Exec(args ...interface{}) (sql.Result, error) {
	var startTime time.Time
	if Enable {
		startTime = time.Now()
	}
	result, error := proxy.Original.Exec(args...)
	if Enable {
		defer proxy.measure(startTime)
	}
	return result, error
}

// instrument Query
func (proxy *Stmt) Query(args ...interface{}) (*sql.Rows, error) {
	var startTime time.Time
	if Enable {
		startTime = time.Now()
	}
	rows, error := proxy.Original.Query(args...)
	if Enable {
		defer proxy.measure(startTime)
	}
	return rows, error
}

// instrument QueryRow
func (proxy *Stmt) QueryRow(args ...interface{}) *sql.Row {
	var startTime time.Time
	if Enable {
		startTime = time.Now()
	}
	row := proxy.Original.QueryRow(args...)
	if Enable {
		defer proxy.measure(startTime)
	}
	return row
}

// just wrap
func (proxy *Stmt) Close() error {
	return proxy.Original.Close()
}
