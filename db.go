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
	if Enable {
		startTime := time.Now()
		defer proxy.measure(startTime, query)
	}
	return proxy.Original.Exec(query, args...)
}

// instrument Query
func (proxy *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if Enable {
		startTime := time.Now()
		defer proxy.measure(startTime, query)
	}
	return proxy.Original.Query(query, args...)
}

// instrument QueryRow
func (proxy *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	if Enable {
		startTime := time.Now()
		defer proxy.measure(startTime, query)
	}
	return proxy.Original.QueryRow(query, args...)
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
