# go-sql\_metrics

Instrument database/sql queries

# Usage

```go
import (
  "database/sql"
  sql_metrics "github.com/sonots/go-sql_metrics"
)

func main() {
  rawDB, _ := sql.Open("sqlite3", "foo.db")
  db := sql_metrics.WrapDB("foo", rawDB) // wrap *sql.DB by *sql_metrics.DB

  // instrument Exec
  result, err := db.Exec("INSERT INTO memos (id, body) VALUES (?)", id, body)
  // instrument Query
  rows, err := db.Query("SELECT body FROM memos WHERE id = ?", id)
  // instrument QueryRow
  err := db.QueryRow("SELECT body FROM memos WHERE id = ?", id).Scan(&body)

  sql_metrics.Verbose = true // print metrics on each query
  sql_metrics.Print(1) // print metrics on each 1 second
  // sql_metrics.Enable = false // turn off the instrumentation
}
```

Output Example (LTSV format):

```
time:2014-09-13 17:35:49.764851359 +0900 JST    db:foo  query:SELECT body FROM memos WHERE id = ?   count:4 max:0.001562    mean:0.000817   min:0.000033    percentile95:0.001562  duration:1
```

Verbose Output Example (LTSV format):

```
time:2014-09-13 17:35:49.717393256 +0900 JST    db:foo  query:SELECT body FROM memos WHERE id = ?   elapsed:0.000910
```

# ToDo

* Write tests

# Contribution

* Fork (https://github.com/sonots/go-sql_metrics/fork)
* Create a feature branch
* Commit your changes
* Rebase your local changes against the master branch
* Create new Pull Request

# Copyright

* See [LICENSE](./LICENSE)
