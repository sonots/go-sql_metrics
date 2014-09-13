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
	result, error := proxy.Original.Exec(query, args...)
	return result, error
}

func (proxy *Tx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	proxy.query = query
	rows, error := proxy.Original.Query(query, args...)
	return rows, error
}

func (proxy *Tx) QueryRow(query string, args ...interface{}) *sql.Row {
	proxy.query = query
	row := proxy.Original.QueryRow(query, args...)
	return row
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
	var startTime time.Time
	if Enable {
		startTime = time.Now()
	}
	if Enable {
		defer proxy.measure(startTime)
	}
	return proxy.Original.Commit()
}

// just wrap
func (proxy *Tx) Rollback() error {
	return proxy.Original.Rollback()
}
