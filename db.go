package sql_metrics

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

type DB struct {
	Original *sql.DB
	*Metrics
}

func newDB(name string, db *sql.DB) *DB {
	return &DB{
		db,
		newMetrics(name),
	}
}

// instrument Exec
func (proxy *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	var startTime time.Time
	if Enable {
		startTime = time.Now()
	}
	result, error := proxy.Original.Exec(query, args...)
	if Enable {
		defer proxy.measure(startTime, query)
	}
	return result, error
}

// instrument Query
func (proxy *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	var startTime time.Time
	if Enable {
		startTime = time.Now()
	}
	rows, error := proxy.Original.Query(query, args...)
	if Enable {
		defer proxy.measure(startTime, query)
	}
	return rows, error
}

// instrument QueryRow
func (proxy *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	var startTime time.Time
	if Enable {
		startTime = time.Now()
	}
	row := proxy.Original.QueryRow(query, args...)
	if Enable {
		defer proxy.measure(startTime, query)
	}
	return row
}

// forward to Stmt
func (proxy *DB) Prepare(query string) (*Stmt, error) {
	stmt, err := proxy.Original.Prepare(query)
	proxyStmt := newStmt(proxy, stmt, query)
	return proxyStmt, err
}

// forward to Tx
func (proxy *DB) Begin() (*Tx, error) {
	tx, err := proxy.Original.Begin()
	proxyTx := newTx(proxy, tx)
	return proxyTx, err
}

func (proxy *DB) Close() error {
	return proxy.Original.Close()
}

func (proxy *DB) Driver() driver.Driver {
	return proxy.Original.Driver()
}

func (proxy *DB) Ping() error {
	return proxy.Original.Ping()
}

func (proxy *DB) SetMaxIdleConns(n int) {
	proxy.Original.SetMaxIdleConns(n)
}

func (proxy *DB) SetMaxOpenConns(n int) {
	proxy.Original.SetMaxOpenConns(n)
}
