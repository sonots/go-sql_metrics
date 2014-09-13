package sql_metrics

import (
	"database/sql"
	"time"
)

type Tx struct {
	proxyDB  *DB
	Original *sql.Tx
	query    string
}

func newTx(proxy *DB, tx *sql.Tx) *Tx {
	return &Tx{
		proxyDB:  proxy,
		Original: tx,
	}
}

func (proxy *Tx) measure(startTime time.Time) {
	proxy.proxyDB.measure(startTime, proxy.query)
}

// instrument Exec
func (proxy *Tx) Exec(query string, args ...interface{}) (sql.Result, error) {
	proxy.query = query
	return proxy.Original.Exec(query, args...)
}

func (proxy *Tx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	proxy.query = query
	return proxy.Original.Query(query, args...)
}

func (proxy *Tx) QueryRow(query string, args ...interface{}) *sql.Row {
	proxy.query = query
	return proxy.Original.QueryRow(query, args...)
}

func (proxy *Tx) Prepare(query string) (*sql.Stmt, error) {
	proxy.query = query
	return proxy.Original.Prepare(query)
}

func (proxy *Tx) Stmt(stmt *sql.Stmt) *sql.Stmt {
	// proxy.query = hmm, how to get query string from sql.Stmt?
	return proxy.Original.Stmt(stmt)
}

// meassure time at commit
func (proxy *Tx) Commit() error {
	if Enable {
		startTime := time.Now()
		defer proxy.measure(startTime)
	}
	return proxy.Original.Commit()
}

// just wrap
func (proxy *Tx) Rollback() error {
	return proxy.Original.Rollback()
}
